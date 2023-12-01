package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"

	recorder "github.com/kajikentaro/request-record-middleware"
	"github.com/kajikentaro/request-record-middleware/types"
)

func handler(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	fmt.Println("handler:", len(b))
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	err := r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rec := recorder.New(types.Setting{OutputDir: "/tmp/hoge"})

	http.Handle("/", rec.Middleware(http.HandlerFunc(handler)))
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
