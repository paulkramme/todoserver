package main

import "fmt"
import "net/http"

type Site struct {
	Title string
	Body []byte
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	fmt.Println("TODO SERVER\nCopyright by Paul Kramme 2017")
	http.HandleFunc("/api", handler)
	http.ListenAndServe(":8080", nil)
}
