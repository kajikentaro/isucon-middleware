//go:build unix

package storages

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/kajikentaro/isucon-middleware/isumid/models"
	"github.com/kajikentaro/isucon-middleware/isumid/settings"
)

// Headers similar to what an ISUCON app behind nginx receives / returns.
var benchReqHeader = map[string][]string{
	"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
	"Accept-Encoding": {"gzip"},
	"Connection":      {"keep-alive"},
	"Content-Type":    {"application/x-www-form-urlencoded"},
	"Cookie":          {"session=MTcwMDAwMDAwMHxEdi1CQkFFQ180SUFBUkFCRUFBQUpQLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FHYzNSeWFXNW5EQVlBQkRFeU16UT18abcdefghijklmnopqrstuvwxyz0123456789"},
	"User-Agent":      {"isucandar"},
}

var benchResHeader = map[string][]string{
	"Content-Type": {"application/json; charset=UTF-8"},
	"Set-Cookie":   {"session=MTcwMDAwMDAwMHxEdi1CQkFFQ180SUFBUkFCRUFBQUpQLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FHYzNSeWFXNW5EQVlBQkRFeU16UT18; Path=/; Expires=Thu, 01 Jan 2099 00:00:00 GMT; Max-Age=86400; HttpOnly"},
	"Vary":         {"Accept-Encoding"},
}

func benchInput() models.RecordedDataInput {
	return models.RecordedDataInput{
		Method:     "GET",
		Url:        "/message?channel_id=1&last_message_id=12345",
		ReqHeader:  benchReqHeader,
		ReqBody:    nil,
		StatusCode: 200,
		ResHeader:  benchResHeader,
		ResBody:    []byte(`[{"id":1,"user":{"name":"foo","display_name":"Foo","avatar_icon":"default.png"},"date":"2017/01/01 00:00:00","content":"hello"}]`),
		Path:       "/message",
		StartedAt:  time.Now(),
		Duration:   1234 * time.Microsecond,
	}
}

func benchMeta() models.Meta {
	in := benchInput()
	return models.Meta{
		Method:     in.Method,
		Url:        in.Url,
		ReqHeader:  in.ReqHeader,
		StatusCode: in.StatusCode,
		ResHeader:  in.ResHeader,
		IsReqText:  false,
		IsResText:  true,
		Ulid:       genUlidStr(),
		ReqLength:  0,
		ResLength:  len(in.ResBody),

		Path:        in.Path,
		StartedAtUs: in.StartedAt.UnixMicro(),
		DurationUs:  in.Duration.Microseconds(),
	}
}

// cpuTime returns user+sys CPU time consumed by this process so far.
func cpuTime(b *testing.B) time.Duration {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		b.Fatal(err)
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

func BenchmarkGenUlid(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		genUlidStr()
	}
}

func BenchmarkGenUlidParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			genUlidStr()
		}
	})
}

// BenchmarkSerializeMeta measures the work done on the request goroutine
// to turn metadata into a row (header encoding etc.).
func BenchmarkSerializeMeta(b *testing.B) {
	meta := benchMeta()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := serializeMeta(meta); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSaveParallel measures Save on the request goroutine: ulid + IsText + serialize + enqueue + body files.
// The background bulk insert is not included (see BenchmarkInsertMetaBulk).
func BenchmarkSaveParallel(b *testing.B) {
	dir := b.TempDir()
	s, err := New(settings.Setting{OutputDir: dir})
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()
	in := benchInput()
	b.ReportAllocs()
	b.ResetTimer()
	cpu := cpuTime(b)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := s.Save(in); err != nil {
				b.Fatal(err)
			}
		}
	})
	s.Flush()
	b.StopTimer()
	b.ReportMetric(float64((cpuTime(b)-cpu).Nanoseconds())/float64(b.N), "cpu-ns/op")
	time.Sleep(100 * time.Millisecond) // let in-flight bulk inserts finish before Close
}

// BenchmarkInsertMetaBulk measures the background sqlite insert, per row.
func BenchmarkInsertMetaBulk(b *testing.B) {
	dir := b.TempDir()
	s, err := New(settings.Setting{OutputDir: dir})
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	const chunk = 50
	rows := make([][]interface{}, 0)
	for i := 0; i < b.N; i += chunk {
		batch := make([]interface{}, 0, chunk)
		for j := 0; j < chunk; j++ {
			m := benchMeta()
			sm, err := serializeMeta(m)
			if err != nil {
				b.Fatal(err)
			}
			batch = append(batch, sm)
		}
		rows = append(rows, batch)
	}
	b.ReportAllocs()
	b.ResetTimer()
	cpu := cpuTime(b)
	for _, batch := range rows {
		insertMetaBulk(s.db, batch)
	}
	b.StopTimer()
	b.ReportMetric(float64((cpuTime(b)-cpu).Nanoseconds())/float64(len(rows)*chunk), "cpu-ns/row")
	os.RemoveAll(dir)
}
