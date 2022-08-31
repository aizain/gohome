# 第一周 web服务支持优雅退出

## TODO 
- Option 设计模式

## 一、概述

启动 web 服务，监听 8080/8081

- 8080 外部服务
- 8081 内部服务

1 使用本地缓存

- 使用 write-back 模式
  - key 过期时刷入数据库
  - 关闭应用时，全部刷入数据库

2 服务优雅退出

- ctrl+c 时，立即拒绝新请求
- 需要等待已有请求处理完毕
- 关闭 8080/8081 两个监听
- 支持注册退出回调，处理本地缓存刷库

```go
// 代表应用本身
type App struct {
	servers []*Server
}

// 代表一个http服务器，一个服务器监听一个端口
type Server struct {
	
}
```


## 二、设计文档

### 1 背景

```
优雅退出：
为了保障程序退出时，能确保正在处理的请求能正常结束，
同时能在应用结束前释放所有的资源，并保障业务的正确性，
不会造成数据丢失。

结束信号：
ctrl+c
```

### 2 名词解释


| 名词     | 含义                        | 注释                                                                                                   |
|--------|---------------------------|------------------------------------------------------------------------------------------------------|
| 应用     | 指部署到特定机器的服务实例，应用和实例具有相同含义 |                                                                                                      |
| server | 监听了某个端口，提供了对外服务的某个类型的实例   | 一般应用和 server 是一对多的关系，在一个应用内会启动多个 server，监听不同的端口；<br/>当前示例应用会启动两个 server，一个对外提供服务，另一个对内提供管理接口，供开发人员使用 |

### 3 需求分析
#### 3.1 场景分析
#### 3.2 功能性需求
#### 3.3 非功能性需求

```
# 优雅退出流程
拒绝新请求=》等待已有请求=》关闭server=》执行回调=》释放资源

# PS
# 1 开发者在回调时可能会使用资源，并且不希望开发者继续使用 server，
# 所以回调在关闭 server 与释放资源之间
# 2 回调上并无排序、优先级、依赖等需求，后续设计将使用并发处理方案，
# 开发者要自行确认回调之前无依赖、顺序
```

### 4 详细设计
#### 4.1 拒绝新请求

可以考虑装饰 http.Server，在每次接收到新请求时，检查是否要拒绝

示例代码：此处代码已写好，可以感受一下装饰器
```go
type serverMux struct {
	reject bool
	*http.ServeMux
}
func (s *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if s.reject {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("服务已关闭"))
		return
    }
	s.ServeMux.ServeHTTP(w, r)
}
```

当开始优雅退出时，reject 标记为 true，如果此时收到新请求，将会返回 503
503 Service Unavailable

#### 4.2 等待已有请求执行完毕

两种思路：
1. 等待一段固定时间
2. 实时维护正在执行的请求数量
当前实现选用等待一段时间，开发者可以配置等待时间

#### 4.3 自定义选项和注册回调

在一些步骤中，开发者希望自定义开发一些参数，比如等待时间。
需要允许开发者注册退出回调函数，可以通过 Option 设计模式来解决。
回调也是需要考虑超时问题的，既不希望回调长时间运行，也希望开发者能明确意识到超时这个问题，
所以采用了 context 的方案，让开发者自己处理超时。
注册回调将会被并发调用执行。

```go
type Option func(*App)
type ShutdownCallback func(ctx context.Context)
func WithShutdownCallbacks func(cbs ...ShutdownCallback) Option {
	return func(app *App) {
		app.cbs = cbs
    }
}
type App struct {
	cbs []ShutdownCallback
}
```

#### 4.4 监听系统信号

GO 监听系统信号是比较简单的，代码如下：

```go
c := make(chan os.Signal, 1)
signal.Notify(c, signals)
select {
    case <-c:
		// 监听到了关闭信号
}
```

不同系统要监听的 signals 是不一样的，可以利用 Go 的编译器特性，为不同平台进行定义，此处使用后缀区分。

_darwin.go
Darwin是由苹果公司于2000年所发布的一个开放源代码操作系统。Darwin是macOS和iOS操作环境的操作系统部分。苹果公司于2000年把Darwin发布给开放源代码社群。
```go

```

_windows.go

_linux.go


#### 4.5 强制退出

采用两次监听的策略，第一次信号采用优雅退出。之后需要两件事情：
- 再次收到退出信号好，立即强制退出
- 启动超时计时器，超时后强制退出