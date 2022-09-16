package main

import (
	"log"
	"reflect"
)

type User struct {
	Acct
	Name  string
	Age   int64
	State int8
	Money float64
	Addr  *Addr
	a     Acct
}

type Acct struct {
	ID string
}

type Addr struct {
	Desc string
}

func main() {
	typeAndValueDemo()
}

func typeAndValueDemo() {
	user := &User{
		Name:  "zain",
		Age:   18,
		State: 1,
		Money: 99.9999,
		Addr: &Addr{
			Desc: "北京市",
		},
		Acct: Acct{
			ID: "1",
		},
		a: Acct{
			ID: "10",
		},
	}
	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)
	log.Printf("user name, %v", t.String())
	log.Printf("user type, %v", t)
	log.Printf("user value, %v", v)
	log.Printf("user actual value, %v", v.Interface())

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	log.Printf("user name, %v", v.Type().Name())
	log.Printf("user type, %v", v)
	log.Printf("user value, %v", v)
	log.Printf("user actual value, %v", v.Interface())

	fCnt := v.NumField()
	for i := 0; i < fCnt; i++ {
		cv := v.Field(i)
		ct := t.Field(i)
		log.Printf("user index:field:type:value, %v:%v:%v:%v", i, ct.Name, ct, cv)
	}

}
