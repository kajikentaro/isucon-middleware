package isumid

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kajikentaro/isucon-middleware/isumid/handlers"
	"github.com/kajikentaro/isucon-middleware/isumid/middlewares"
	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/services"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
)

type Recorder struct {
	// Middleware func(http.Handler) http.Handler
	handler    handlers.Handler
	middleware middlewares.Middleware
}

func (rec Recorder) Middleware(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/fetch-req-body", rec.handler.FetchAll)
	mux.HandleFunc("/fetch-res-body", rec.handler.FetchAll)
	mux.HandleFunc("/fetch-reproduced-res-body", rec.handler.FetchAll)
	mux.HandleFunc("/fetch-all", rec.handler.FetchAll)
	mux.Handle("/reproduce", rec.middleware.Reproducer(next))
	mux.Handle("/", rec.middleware.Recorder(next))
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

	ser := services.New(storage)
	han := handlers.New(ser)

	mid := middlewares.New(storage)

	return Recorder{handler: han, middleware: mid}
}
