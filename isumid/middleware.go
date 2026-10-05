package isumid

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kajikentaro/isucon-middleware/isumid/handlers"
	"github.com/kajikentaro/isucon-middleware/isumid/middlewares"
	"github.com/kajikentaro/isucon-middleware/isumid/services"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/kajikentaro/isucon-middleware/isumid/storages"
)

type AutoSwitch = settings.AutoSwitch
type Setting = settings.Setting

type Recorder struct {
	// Middleware func(http.Handler) http.Handler
	handler    handlers.Handler
	middleware middlewares.Middleware
	prefix     string
}

func (rec *Recorder) Middleware(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	// Routes are registered with the prefix so that ServeMux's trailing-slash redirects
	// (e.g. "/isumid" -> "/isumid/", "/isumid/req-body" -> "/isumid/req-body/") keep the prefix.
	// The handlers receive the path with the prefix stripped.
	handle := func(pattern string, handler http.Handler) {
		mux.Handle(rec.prefix+pattern, http.StripPrefix(rec.prefix, handler))
	}
	handle("/start-recording", http.HandlerFunc(rec.middleware.StartRecording))
	handle("/stop-recording", http.HandlerFunc(rec.middleware.StopRecording))
	handle("/is-recording", http.HandlerFunc(rec.middleware.IsRecording))
	handle("/req-body/", http.HandlerFunc(rec.handler.FetchReqBody))
	handle("/res-body/", http.HandlerFunc(rec.handler.FetchResBody))
	handle("/remove/", http.HandlerFunc(rec.handler.Remove))
	handle("/remove-all", http.HandlerFunc(rec.handler.RemoveAll))
	handle("/reproduced-res-body/", http.HandlerFunc(rec.handler.FetchReproducedResBody))
	handle("/search", http.HandlerFunc(rec.handler.Search))
	handle("/export", http.HandlerFunc(rec.handler.Export))
	handle("/export-size", http.HandlerFunc(rec.handler.ExportSize))
	handle("/reproduce/", rec.middleware.Reproducer(next))
	handle("/", rec.handler.Frontend())
	mux.Handle("/", rec.middleware.Recorder(next))
	return mux
}

func New(options *Setting) *Recorder {
	defaultSetting := Setting{
		Prefix:        "/isumid",
		OutputDir:     filepath.Join(os.TempDir(), "isumid"),
		RecordOnStart: false,
		AutoStop:      nil,
		AutoStart:     nil,
	}

	if options == nil {
		options = &defaultSetting
	} else {
		if options.Prefix == "" {
			options.Prefix = defaultSetting.Prefix
		}
		if !strings.HasPrefix(options.Prefix, "/") {
			options.Prefix = "/" + options.Prefix
		}
		if options.OutputDir == "" {
			options.OutputDir = defaultSetting.OutputDir
		}
	}

	// DI
	storage, err := storages.New(*options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	ser := services.New(storage)
	han := handlers.New(ser)

	mid := middlewares.New(storage, options)

	return &Recorder{handler: han, middleware: mid, prefix: options.Prefix}
}
