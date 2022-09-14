package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	taskPoolDemo(10)

}

func taskPoolDemo(num int) {
	p := sync.Pool{}
	for i := 0; i < num; i++ {
		p.Put(i)
	}

	w := sync.WaitGroup{}
	for i := 0; i < num*10; i++ {
		w.Add(1)
		go func(i int) {
			e := p.Get()
			for e == nil {
				time.Sleep(time.Second)
				e = p.Get()
			}
			fmt.Printf("get from pool, i,e: %v,%v\n", i, e)
			p.Put(e)
			w.Done()
		}(i)
	}
	w.Wait()

}
