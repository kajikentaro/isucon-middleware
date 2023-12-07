package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"

	"github.com/kajikentaro/isucon-middleware/isumid"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/labstack/echo/v4"
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
	echo.New()
	rec := isumid.New(models.Setting{OutputDir: "/tmp/hoge"})

	http.Handle("/", rec.Middleware(http.HandlerFunc(handler)))
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
