# 第二周 ⽣成 INSERT 语句

## 一、概述

利用反射生成能在MySQL上执行的语句，已有预先定义好的方法

```go
package xxx

func InsertStmt(entity any) (string, []any, error) { }
```

目前只考虑一下几种输入情况：
- nil
- 结构体
- 一级指针，指针指向的是结构体
- 组合，但是不能是指针形态的组合


数据库驱动只考虑以下类型：
- 基本类型
- string, []byte, time.Time
- 实现了 driver.Valuer / sql.Scanner 两个接口的类型


处理流程：
- 普通字段直接处理，它的名作为列名，它的值作为参数
- 组合需要递归遍历该结构体


## 二、设计文档

### 1 背景

```
在实际业务中会遇到需要将一个数据实例入库的情况，这时候就面临如何将实例转化为SQL语句的问题，并且还需要生成其参数

以下情况需要考虑
批量插入
指定列，以及指定列的表达式
组合：有很多公司的数据实例包含一些公共结构体，通过组合方式嵌入
upsert：数据不存在就插入，存在就删除

```

### 2 名称解释

多重组合：一个结构体同时组合了多个结构体
深层组合：一个结构体 A 组合了另一个结构体 B，而 B 本身还结合了 C
多级指针：指向指针的指针，如 **int
方言：不同数据库对SQL的支持包括语法都有一些差异，不同数据库支持的 SQL 认为是方言
upsert：如果不存在就插入，如果数据存在就更新。通过主键/唯一索引进行判断

### 3 需求分析

#### 3.1 场景分析

主要是考虑以下因素：
1. 开发者如何构建输入，核心是使用指针还是非指针
2. 开发者如何定义模型，核心是是否使用组合以及如何组合
3. 开发者定义何种字段，核心是复杂结构体作为了字段
4. 批量插入
5. 指定列，及列的表达式
6. 使用 upsert （指定唯一列，指定更新的列及值）

本次需要支持的场景

| 场景          | 例子                         | 说明             |
|-------------|----------------------------|----------------|
| 使用结构体作为输入   | InsertStmt(User{})         | 能够正常处理         |
| 使用指针作为输入    | InsertStmt(&User{})        | 能够正常处理         |
| 使用 nil 作为输入 | InsertStmt(nil)            | 不需要支持，返回错误     |
| 使用了组合       | type Buyer struct { User } | 支持，要包含 User 字段 |

#### 3.2 方言分析

MySQL INSERT 主要是三种形态

```sql
-- 具有 ON DUPLICATE KEY UPDATE 解决 upsert 问题
INSERT [LOW_PRIORITY | DELAYED | HIGH_PRIORITY] [IGNORE]
[INTO] tbl_name
[PARTITION (partition_name [, partition_name] ...)]
[(col_name [, col_name] ...)]
{ {VALUES | VALUE} (value_list) [, (value_list)] ... }
[AS row_alias[(col_alias [, col_alias] ...)]]
[ON DUPLICATE KEY UPDATE assignment_list]
```

```sql
-- 使用了 SET assignment_list
INSERT [LOW_PRIORITY | DELAYED | HIGH_PRIORITY] [IGNORE]
[INTO] tbl_name
[PARTITION (partition_name [, partition_name] ...)]
SET assignment_list
[AS row_alias[(col_alias [, col_alias] ...)]]
[ON DUPLICATE KEY UPDATE assignment_list]
```

```sql
-- 使用了 SELECT 子句
INSERT [LOW_PRIORITY | HIGH_PRIORITY] [IGNORE]
[INTO] tbl_name
[PARTITION (partition_name [, partition_name] ...)]
[(col_name [, col_name] ...)]
{ SELECT ...
| TABLE table_name
| VALUES row_constructor_list
}
[ON DUPLICATE KEY UPDATE assignment_list]
```

assignment_list 定义

```sql
assignment:
col_name =
value
| [row_alias.]col_name
| [tbl_name.]col_name
| [row_alias.]col_alias
assignment_list:
assignment [, assignment] ...
```

关键是 assignment 有四种：
1. value：⼀个纯粹的值
2. row_alias.col_name：在使⽤了⾏别名的时候才会有的形态
3. tbl_name.col_name：指定更新为插⼊的值
4. row_alias.col_alias：这个和第⼆种⽐较像，所不同的是这⾥使⽤的是别名

#### 3.3 功能性需求

##### 3.3.1 生成查询

将结构体转化为 INSERT 查询语句
- SQL
- 参数

要支持以下选项
- 单个插入
- upsert

#### 3.4 非功能性需求

- 扩展性。针对未来更复杂


### 4 设计

#### 4.1 总体设计

整体采用 Builder 模式进行，复杂逻辑可以拆为单一方法，使用泛型避免插入不同类型的数据

```go
package xxx

type Inserter[T any] struct {
}

// 构建 SQL
func (i *Inserter[T]) Build() (Query, error) {}

type Query struct {
    SQL string
    Args []any
}
```

#### 4.2 详细设计

Insert

```go
package xxx

func (i *Inserter[T]) Values(vals...*T) *Inserter {
    panic("implement me")
}

func (i *Inserter[T]) Columns(cols...string) *Inserter[T] {
	panic("implement me")
}

func (i *Inserter[T]) Exec(ctx context.Context) (sql.Result, error) {
	panic("implement me")
}
```

Upsert

```go
package xxx
type Upsert struct {
	doNothing bool
	updateColumns []string
	conflictColumns []string
}
type UpsertBuilder[T any] struct {
	i *Inserter[T]
	conflictColumns []string
}
// Update 指定在冲突的时候要执行 update 的列
func (u *UpsertBuilder[T]) Update(cols...string) *Inserter[T] {
	i.upsert = Upsert{conflictColumns: u.conflictColumns, updateColumns: cols}
    return i
}

func (u *UpsertBuilder[T]) DoNothing() *Inserter[T] {
    i.upsert = Upsert{conflictColumns: u.conflictColumns, doNothing: true}
}
func (u *UpsertBuilder[T]) ConflictColumns(cols...string) *UpsertBuilder[T]
    u.conflictColumns = cols
}

func (i *Inserter[T]) Upsert() *UpsertBuilder[T] {
	return &UpsertBuilder {
		i: i
	}
}
```

#### 4.3 例子

```go
type User struct {
    ID uint64
    Email string
    FirstName string
    Age uint8
}

NewInserter[User]().Values(&User{ID: 1, Email: "xxx@xx"})
// INSERT INTO `user`(id, email, first_name, age) VALUES(?,?,?,?)
// []{1, "xxx@xx", "", 0}
```

### 5 测试

#### 5.1 单元测试


