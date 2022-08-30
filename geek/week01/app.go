package main

import (
	"fmt"
	"net/http"
	"time"
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("request path: %v\n", request.URL.Path)
	fmt.Printf("request host: %v\n", request.URL.Host)
	fmt.Printf("request query: %v\n", request.URL.RawQuery)
	fmt.Printf("request host: %v\n", request.Host)

	fmt.Fprintf(writer, "hello 23333 %v\n", request.URL.Path)
}

func workHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("begin work\n")
	time.Sleep(5 * time.Second)
	fmt.Printf("process work\n")
	time.Sleep(5 * time.Second)
	fmt.Printf("end work\n")

	fmt.Fprintf(writer, "1000")
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/work", workHandler)
	http.ListenAndServe(":8080", nil)

}
