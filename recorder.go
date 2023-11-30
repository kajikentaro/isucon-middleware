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
)

type readCloser struct {
	io.Reader
	close func() error
}

func (n readCloser) Close() error {
	return n.close()
}

func (rec Recorder) HandlerFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		reader := io.TeeReader(r.Body, &buf)
		newreadcloser := readCloser{
			Reader: reader,
			close:  r.Body.Close,
		}
		r.Body = newreadcloser

		next.ServeHTTP(w, r)

		body, err := io.ReadAll(&buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println(len(body))

		ulid, err := ulid.New(ulid.Timestamp(time.Now()), nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		outPath := filepath.Join(rec.OutputDir, ulid.String())
		err = os.WriteFile(outPath, body, 0666)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})
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
		OutputDir: "/tmp/request-record-middleware",
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
