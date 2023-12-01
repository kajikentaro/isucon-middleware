package recorders

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kajikentaro/request-record-middleware/types"
	"github.com/oklog/ulid"
	"github.com/vmihailenco/msgpack/v5"
)

type Recorder struct {
	setting types.Setting
}

func New(setting types.Setting) Recorder {
	return Recorder{setting: setting}
}

func (r Recorder) Middleware(header http.Header, body io.Reader) {
	fmt.Println("recorder")

	ulid, err := ulid.New(ulid.Timestamp(time.Now()), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	ulidStr := ulid.String()

	outPathBody := filepath.Join(r.setting.OutputDir, ulidStr+".body")

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err = os.WriteFile(outPathBody, bodyBytes, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	outPathHeader := filepath.Join(r.setting.OutputDir, ulidStr+".header")
	headerBytes, err := msgpack.Marshal(header)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	err = os.WriteFile(outPathHeader, headerBytes, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
