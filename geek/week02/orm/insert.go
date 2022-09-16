package orm

import (
	"fmt"
	"reflect"
	"strings"
)

type Query struct {
	SQL  string
	Args []any
}

func NewInserter[T any]() *Inserter[T] {
	return &Inserter[T]{}
}

type Inserter[T any] struct {
	values []*T
}

func (i *Inserter[T]) Build() (*Query, error) {
	return i.buildByReflect()
}

func (i *Inserter[T]) Values(values ...*T) *Inserter[T] {
	i.values = values
	return i
}

func (i *Inserter[T]) buildByReflect() (*Query, error) {
	values := make([]any, 0, len(i.values)*10)
	sqls := make([]string, 0, len(i.values))
	for _, value := range i.values {
		// INSERT INTO `user`(id, email, first_name, age) VALUES(?,?,?,?)
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		fields := make([]string, 0, v.NumField())
		placeholders := make([]string, 0, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.Kind() == reflect.Struct {
				for i := 0; i < field.NumField(); i++ {
					innerField := field.Field(i)
					fieldStruct := field.Type().Field(i)
					fields = append(fields, strings.ToLower(fieldStruct.Name))
					placeholders = append(placeholders, "?")
					values = append(values, innerField.Interface())
				}
			} else {
				fieldStruct := v.Type().Field(i)
				fields = append(fields, strings.ToLower(fieldStruct.Name))
				placeholders = append(placeholders, "?")
				values = append(values, field.Interface())
			}
		}
		sqls = append(sqls, fmt.Sprintf("INSERT INTO `%v`(%v) VALUES(%v)",
			strings.ToLower(v.Type().Name()), strings.Join(fields, ","), strings.Join(placeholders, ",")))
	}
	return &Query{
		SQL:  strings.Join(sqls, ";"),
		Args: values,
	}, nil
}
