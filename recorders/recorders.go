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

func (r Recorder) Middleware(reqHeader http.Header, reqBody io.Reader, resHeader http.Header, resBody *[]byte) {
	// generate ulid
	ulid, err := ulid.New(ulid.Timestamp(time.Now()), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	ulidStr := ulid.String()

	// save request body data
	{
		path := filepath.Join(r.setting.OutputDir, ulidStr+".req.body")
		data, err := io.ReadAll(reqBody)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	// save request header data
	{
		path := filepath.Join(r.setting.OutputDir, ulidStr+".req.header")
		data, err := msgpack.Marshal(reqHeader)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	// save response body data
	{
		path := filepath.Join(r.setting.OutputDir, ulidStr+".res.body")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = os.WriteFile(path, *resBody, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	// save response header data
	{
		path := filepath.Join(r.setting.OutputDir, ulidStr+".req.header")
		data, err := msgpack.Marshal(resHeader)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
