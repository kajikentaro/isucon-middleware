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
	fmt.Println("DATA handler:", len(b))
	fmt.Println("handler end")
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	err := r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	hoge := recorder.New(recorder.RecorderOptions{OutputDir: "/tmp/hoge"})

	http.Handle("/", hoge.HandlerFunc(http.HandlerFunc(handler)))
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
