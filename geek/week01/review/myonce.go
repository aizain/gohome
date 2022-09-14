package main

import (
	"fmt"
	"sync"
)

//func main() {
//	//initDemo()
//	instanceDemo()
//}

var instance *res
var instanceOnce sync.Once

func instanceDemo() {
	instanceOnce.Do(func() {
		instance = &res{}
		instance.Init()
		instance.Init()
	})
	fmt.Printf("instance demo %v\n", instance)
}

type res struct {
	once sync.Once
}

func (r *res) Init() {
	r.once.Do(func() {
		fmt.Printf("初始化执行一次\n")
	})
}

func initDemo() {
	r := res{}
	r.Init()
	r.Init()
	r.Init()
	fmt.Printf("init demo %v\n", r)
}
