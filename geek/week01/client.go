package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func main() {
	doWork()
}

func doWork() {
	resp, err := http.Get("http://localhost:8080/work")
	if err != nil {
		panic("work failed")
	}
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic("work but no money")
	}
	fmt.Printf("money: %v", string(body))
}

func doRootHello() {
	msg := []byte("hello")
	//var msg []byte
	fmt.Printf("msg: %v\n", reflect.TypeOf(msg))

	resp, err := http.Post("http://localhost:8080", "", bytes.NewBuffer(msg))
	if err != nil {
		panic("oh no !")
	} else {
		fmt.Printf("ok\n")
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic("no body no no no")
	}

	fmt.Printf("body: %v\n", string(body))
}
