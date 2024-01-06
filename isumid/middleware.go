package isumid

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kajikentaro/isucon-middleware/isumid/handlers"
	"github.com/kajikentaro/isucon-middleware/isumid/middlewares"
	"github.com/kajikentaro/isucon-middleware/isumid/services"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
)

type Recorder struct {
	// Middleware func(http.Handler) http.Handler
	handler    handlers.Handler
	middleware middlewares.Middleware
}

func (rec Recorder) Middleware(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/isumid/start-recording", rec.middleware.StartRecording)
	mux.HandleFunc("/isumid/stop-recording", rec.middleware.StopRecording)
	mux.HandleFunc("/isumid/is-recording", rec.middleware.IsRecording)
	mux.HandleFunc("/isumid/req-body/", rec.handler.FetchReqBody)
	mux.HandleFunc("/isumid/res-body/", rec.handler.FetchResBody)
	mux.HandleFunc("/isumid/reproduced-res-body/", rec.handler.FetchReproducedResBody)
	mux.HandleFunc("/isumid/list", rec.handler.FetchList)
	mux.Handle("/isumid/reproduce/", rec.middleware.Reproducer(next))
	mux.HandleFunc("/isumid/", rec.handler.Frontend)
	mux.Handle("/", rec.middleware.Recorder(next))
	return mux
}

func New(options *settings.Setting) Recorder {
	def := settings.Setting{
		OutputDir:     filepath.Join(os.TempDir(), "isumid"),
		RecordOnStart: true,
	}

	if options == nil {
		options = &def
	} else {
		if options.OutputDir == "" {
			options.OutputDir = def.OutputDir
		}
	}

	err := os.MkdirAll(options.OutputDir, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// DI
	storage := storages.New(*options)

	ser := services.New(storage)
	han := handlers.New(ser)

	mid := middlewares.New(storage, options)

	return Recorder{handler: han, middleware: mid}
}
