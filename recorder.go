package recorder

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/oklog/ulid"
	"github.com/vmihailenco/msgpack/v5"
)

type readCloser struct {
	io.Reader
	close func() error
}

func (n readCloser) Close() error {
	return n.close()
}

func (rec Recorder) recorderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		reader := io.TeeReader(r.Body, &buf)
		newreadcloser := readCloser{
			Reader: reader,
			close:  r.Body.Close,
		}
		r.Body = newreadcloser

		next.ServeHTTP(w, r)

		ulid, err := ulid.New(ulid.Timestamp(time.Now()), nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		ulidStr := ulid.String()

		outPathBody := filepath.Join(rec.OutputDir, ulidStr)
		err = os.WriteFile(outPathBody, buf.Bytes(), 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		outPathHeader := filepath.Join(rec.OutputDir, ulidStr+".header")
		header, err := msgpack.Marshal(r.Header)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = os.WriteFile(outPathHeader, header, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})
}

func (rec Recorder) hoge(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOGE"))
}

func (rec Recorder) HandlerFunc(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/get-all", rec.hoge)
	mux.Handle("/", rec.recorderMiddleware(next))
	return mux
}

type Middleware func(http.Handler) http.Handler

type Recorder struct {
	RecorderOptions
}

type RecorderOptions struct {
	OutputDir string
}

func New(options RecorderOptions) Recorder {
	def := RecorderOptions{
		OutputDir: filepath.Join(os.TempDir(), "request-record-middleware"),
	}

	if options.OutputDir == "" {
		options.OutputDir = def.OutputDir
	}

	err := os.MkdirAll(options.OutputDir, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	res := Recorder{RecorderOptions: options}
	return res
}
