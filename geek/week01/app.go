package main

import (
	"fmt"
	"net/http"
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("request path: %v\n", request.URL.Path)
	fmt.Printf("request host: %v\n", request.URL.Host)
	fmt.Printf("request query: %v\n", request.URL.RawQuery)
	fmt.Printf("request host: %v\n", request.Host)

	fmt.Fprintf(writer, "hello 23333 %v\n", request.URL.Path)
}

type Server struct {
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", nil)

}
