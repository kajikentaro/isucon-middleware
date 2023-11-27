package recorder

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Recorder struct {
	Sniffer Middleware
}

type newReadCloser struct {
	io.Reader
	close func() error
}

func (n newReadCloser) Close() error {
	fmt.Println("closed")
	return n.close()
}

func HandlerFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("start")

		var buf bytes.Buffer

		reader := io.TeeReader(r.Body, &buf)
		newreadcloser := newReadCloser{
			Reader: reader,
			close:  r.Body.Close,
		}
		r.Body = newreadcloser

		next.ServeHTTP(w, r)

		fmt.Println("LEN middleware:", buf.Len())
		body, err := io.ReadAll(&buf)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		fmt.Println("DATA middleware:", len(body))

		fmt.Println("end")
	})
}

func New() Recorder {
	var r Recorder
	r.Sniffer = HandlerFunc
	return r
}
