package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"

	recorder "github.com/kajikentaro/request-record-middleware"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler start")
	b, _ := io.ReadAll(r.Body)
	fmt.Println(string(b))
	fmt.Println("handler end")
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	hoge := recorder.New()

	http.Handle("/", hoge.Sniffer(http.HandlerFunc(handler)))
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
