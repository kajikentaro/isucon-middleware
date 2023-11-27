package recorder

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/gommon/log"
)

type Middleware func(http.Handler) http.Handler

type Recorder struct {
	Sniffer Middleware
}

func HandlerFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("start")
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Println(string(b))

		next.ServeHTTP(w, r)
		fmt.Println("end")
	})
}

func New() Recorder {
	var r Recorder
	r.Sniffer = HandlerFunc
	return r
}
