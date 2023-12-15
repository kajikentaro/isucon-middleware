package main

import (
	"fmt"
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
	w.Write(b)

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
