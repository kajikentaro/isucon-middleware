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
	originalClose func() error
}

func (n readCloser) Close() error {
	return n.originalClose()
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

type responseWriter struct {
	original    http.ResponseWriter
	writtenData *[]byte
}

func (r responseWriter) Header() http.Header {
	return r.original.Header()
}
func (r responseWriter) Write(in []byte) (int, error) {
	*r.writtenData = in
	return r.original.Write(in)
}
func (r responseWriter) WriteHeader(statusCode int) {
	r.original.WriteHeader(statusCode)
}

func (h Handler) RecorderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		reader := io.TeeReader(r.Body, &buf)
		sniffedReadCloser := readCloser{
			Reader:        reader,
			originalClose: r.Body.Close,
		}
		r.Body = sniffedReadCloser

		sniffedResponseWriter := responseWriter{original: w, writtenData: &[]byte{}}

		next.ServeHTTP(sniffedResponseWriter, r)

		h.recorder.Middleware(r.Header, &buf, w.Header(), sniffedResponseWriter.writtenData)
	})
}
