package main

import "net/http"

func main() {
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})
	http.ListenAndServe(":2345", nil)
}