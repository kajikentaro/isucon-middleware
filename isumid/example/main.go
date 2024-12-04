// this file is for debugging and testing Isucon Middleware
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kajikentaro/isucon-middleware/isumid"
	"github.com/labstack/echo/v4"
)

func handler(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	fmt.Println("handler:", len(b))
	w.Write([]byte("test response"))
	w.Header().Add("sample header", "sample header")

	err := r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	echo.New()
	rec := isumid.New(nil)

	http.Handle("/", rec.Middleware(http.HandlerFunc(handler)))
	fmt.Println("server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
