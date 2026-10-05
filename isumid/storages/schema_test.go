package storages

import (
	"testing"
	"time"

	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
	"github.com/stretchr/testify/assert"
)

func newTestStorage(t *testing.T) Storage {
	s, err := New(settings.Setting{OutputDir: t.TempDir()})
	assert.NoError(t, err)
	t.Cleanup(func() { s.Close() })
	return s
}

func TestColumnsForAnalysis(t *testing.T) {
	s := newTestStorage(t)

	err := s.Save(models.RecordedDataInput{
		Method: "POST",
		Url:    "/login?next=%2F",
		Path:   "/login",
		ReqHeader: map[string][]string{
			"Cookie":       {"session=old"},
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		ReqBody:    []byte("name=foo"),
		StatusCode: 303,
		ResHeader: map[string][]string{
			"Set-Cookie":       {"session=new; Path=/", "flash=; Max-Age=0"},
			"Content-Type":     {"text/html; charset=UTF-8"},
			"Content-Encoding": {"gzip"},
		},
		StartedAt: time.UnixMicro(1700000000000001),
		Duration:  2345 * time.Microsecond,
	})
	assert.NoError(t, err)
	err = s.Save(models.RecordedDataInput{Method: "GET", Url: "/", Path: "/", ReqHeader: map[string][]string{}, ResHeader: map[string][]string{}})
	assert.NoError(t, err)
	s.Flush()

	type row struct {
		Path        string  `db:"path"`
		StartedAtUs int64   `db:"startedAtUs"`
		DurationUs  int64   `db:"durationUs"`
		Cookie      *string `db:"cookie"`
		SetCookie   *string `db:"setCookie"`
		ResEncoding *string `db:"resEncoding"`
	}
	// headers can be read with the JSON functions of SQLite
	query := `
		SELECT path, startedAtUs, durationUs,
			json_extract(reqHeader, '$.Cookie[0]') AS cookie,
			json_extract(resHeader, '$.Set-Cookie[1]') AS setCookie,
			json_extract(resHeader, '$.Content-Encoding[0]') AS resEncoding
		FROM metadata ORDER BY startedAtUs DESC`
	var rows []row
	assert.NoError(t, s.db.Select(&rows, query))
	assert.Len(t, rows, 2)

	str := func(s string) *string { return &s }
	assert.Equal(t, row{
		Path:        "/login",
		StartedAtUs: 1700000000000001,
		DurationUs:  2345,
		Cookie:      str("session=old"),
		SetCookie:   str("flash=; Max-Age=0"),
		ResEncoding: str("gzip"),
	}, rows[0])

	// absent headers are NULL
	assert.Equal(t, row{Path: "/"}, rows[1])

	// Meta exposes the new fields too
	metaList, _, err := s.SearchMetaList("/login%", 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2345), metaList[0].DurationUs)
}

func TestGenUlidUnique(t *testing.T) {
	const goroutines = 8
	const perGoroutine = 10000
	results := make(chan []string, goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			ids := make([]string, 0, perGoroutine)
			for i := 0; i < perGoroutine; i++ {
				ids = append(ids, genUlidStr())
			}
			results <- ids
		}()
	}
	seen := map[string]bool{}
	for g := 0; g < goroutines; g++ {
		for _, id := range <-results {
			assert.False(t, seen[id], "duplicated ulid: %s", id)
			seen[id] = true
		}
	}
}
