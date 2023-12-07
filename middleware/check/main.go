package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"

	recorder "github.com/kajikentaro/request-record-middleware/middleware"
	"github.com/kajikentaro/request-record-middleware/middleware/models"
)

func handler(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	fmt.Println("handler:", len(b))
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	w.Header().Set("Content-Type", "text/plain")

	err := r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rec := recorder.New(models.Setting{OutputDir: "/tmp/hoge"})

	http.Handle("/", rec.Middleware(http.HandlerFunc(handler)))
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
