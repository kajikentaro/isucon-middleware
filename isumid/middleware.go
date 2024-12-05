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

func (rec *Recorder) Middleware(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/isumid/start-recording", rec.middleware.StartRecording)
	mux.HandleFunc("/isumid/stop-recording", rec.middleware.StopRecording)
	mux.HandleFunc("/isumid/is-recording", rec.middleware.IsRecording)
	mux.HandleFunc("/isumid/req-body/", rec.handler.FetchReqBody)
	mux.HandleFunc("/isumid/res-body/", rec.handler.FetchResBody)
	mux.HandleFunc("/isumid/remove/", rec.handler.Remove)
	mux.HandleFunc("/isumid/remove-all", rec.handler.RemoveAll)
	mux.HandleFunc("/isumid/reproduced-res-body/", rec.handler.FetchReproducedResBody)
	mux.HandleFunc("/isumid/search", rec.handler.Search)
	mux.Handle("/isumid/reproduce/", rec.middleware.Reproducer(next))
	mux.HandleFunc("/isumid/", rec.handler.Frontend)
	mux.Handle("/", rec.middleware.Recorder(next))
	return mux
}

func New(options *settings.Setting) *Recorder {
	defaultSetting := settings.Setting{
		OutputDir:     filepath.Join(os.TempDir(), "isumid"),
		RecordOnStart: false,
		AutoStop:      nil,
		AutoStart:     nil,
	}

	if options == nil {
		options = &defaultSetting
	} else {
		if options.OutputDir == "" {
			options.OutputDir = defaultSetting.OutputDir
		}
	}

	// DI
	storage, err := storages.New(*options)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create a sqlite database", err.Error())
	}
	if err := storage.CreateDir(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to create a directory", err.Error())
		os.Exit(1)
	}

	ser := services.New(storage)
	han := handlers.New(ser)

	mid := middlewares.New(storage, options)

	return &Recorder{handler: han, middleware: mid}
}
