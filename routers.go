package routers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kajikentaro/request-record-middleware/handlers"
	"github.com/kajikentaro/request-record-middleware/models"
	"github.com/kajikentaro/request-record-middleware/services_endpoint"
	"github.com/kajikentaro/request-record-middleware/services_recorder"
	"github.com/kajikentaro/request-record-middleware/storages"
)

type Recorder struct {
	// Middleware func(http.Handler) http.Handler
	handler handlers.Handler
}

func (rec Recorder) Middleware(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/fetch-all", rec.handler.FetchAll)
	mux.Handle("/", rec.handler.RecorderMiddleware(next))
	return mux
}

func New(options models.Setting) Recorder {
	def := models.Setting{
		OutputDir: filepath.Join(os.TempDir(), "request-record-middleware"),
	}

	if options.OutputDir == "" {
		options.OutputDir = def.OutputDir
	}

	err := os.MkdirAll(options.OutputDir, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// DI
	storage := storages.New(options)

	ser := services_endpoint.New(storage)
	rec := services_recorder.New(storage)

	han := handlers.New(ser, rec)
	return Recorder{handler: han}
}
