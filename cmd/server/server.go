package main

import (
	"net/http"
	"log"
	"fmt"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("[error] ", err.Error())
		return
	}
	fmt.Printf("method=%s\n", r.Method)
	fmt.Printf("remoteaddr=%s\n", r.RemoteAddr)
	fmt.Printf("url=%s\n", r.URL)
	for k, arr := range r.Header {
		fmt.Printf("\theaders: %s = %v\n", k, arr)
	}
	println()
}