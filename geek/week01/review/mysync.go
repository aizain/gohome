package main

import (
	"fmt"
	"sync"
	"time"
)

//func main() {
//	mutexDemo()
//}

type SafeMap struct {
	m       map[string]string
	mutex   sync.Mutex
	rwmutex sync.RWMutex
}

func (s *SafeMap) put(key string, val string) string {
	s.mutex.Lock()
	v := s.m[key]
	s.m[key] = val
	s.mutex.Unlock()
	return v
}

func (s *SafeMap) getSet(key string, val string) (string, bool) {
	s.rwmutex.RLock()
	v, ok := s.m[key]
	s.rwmutex.RUnlock()
	if ok {
		return v, true
	}

	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()
	v, ok = s.m[key]
	if ok {
		return v, true
	}
	s.m[key] = val
	return v, false
}

func mutexDemo() {
	s := &SafeMap{
		m: make(map[string]string, 10),
	}
	s.put("a", "v0")

	w := &sync.WaitGroup{}
	w.Add(3)
	go func() {
		time.Sleep(time.Second)
		v := s.put("a", "v1")
		fmt.Printf("map1: %v\n", v)

		w.Done()
	}()
	go func() {
		time.Sleep(time.Second)
		v := s.put("a", "v2")
		fmt.Printf("map2: %v\n", v)

		w.Done()
	}()
	go func() {
		time.Sleep(time.Second)
		v := s.put("a", "v3")
		fmt.Printf("map3: %v\n", v)

		w.Done()
	}()
	w.Wait()
	fmt.Printf("map: %v\n", s.m["a"])
}
