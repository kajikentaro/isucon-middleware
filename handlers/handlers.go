package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/kajikentaro/request-record-middleware/recorders"
	"github.com/kajikentaro/request-record-middleware/services"
)

type Handler struct {
	service  services.Service
	recorder recorders.Recorder
}

type readCloser struct {
	io.Reader
	close func() error
}

func (n readCloser) Close() error {
	return n.close()
}

func New(service services.Service, recorder recorders.Recorder) Handler {
	return Handler{service: service, recorder: recorder}
}

func outputErr(w http.ResponseWriter, err error) {
	message := fmt.Sprintf("%#v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

func (h Handler) FetchAll(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.FetchAll()
	if err != nil {
		outputErr(w, err)
		return
	}
	w.Write([]byte(res))
}

func (h Handler) RecorderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		reader := io.TeeReader(r.Body, &buf)
		newreadcloser := readCloser{
			Reader: reader,
			close:  r.Body.Close,
		}
		r.Body = newreadcloser

		next.ServeHTTP(w, r)

		h.recorder.Middleware(r.Header, &buf)
	})
}
