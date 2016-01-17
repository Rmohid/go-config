// Template for console based debug tracing

package dbg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rmohid/go-config/config"
)

var (
	writers = map[string](io.Writer){}
)

func init() {
	writers[""] = makeDevnull()
	writers["devnull"] = makeDevnull()
	writers["stdout"] = os.Stdout
	writers["stderr"] = os.Stderr
	writers["http"] = makeHttpWriter()
	writers["file"] = makeFileWriter()

	// define all default options
	var opts = [][]string{
		{"dbg.debugWriter", "stderr", "debug log output sink, one of [stderr,stdout,http]"},
		{"dbg.verbosity", "0", "verbosity level for debug output"},
		{"dbg.logfile", "config.log", "filename for log collection"},
		{"dbg.httpUrl", "", "http server for log delivery via http GET"},
	}

	config.PushArgs(opts)
}
func ErrLog(verbosity int, format string, a ...interface{}) {
	val, err := strconv.Atoi(config.Get("dbg.verbosity"))
	if err != nil {
		log.Fatal("dbg.ErrLog:", err)
	}
	if val >= verbosity {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}
func Log(verbosity int, a ...interface{}) {
	val, err := strconv.Atoi(config.Get("dbg.verbosity"))
	if err != nil {
		log.Fatal("dbg.Log:", err)
	}
	if val >= verbosity {
		fmt.Fprintln(writers[config.Get("dbg.debugWriter")],
			a...)
	}
}
func makeDevnull() io.Writer {
	null, err := os.Open(os.DevNull)
	if err != nil {
		log.Fatal("dbg.devnull:", err)
	}
	return null
}

func makeHttpWriter() io.Writer {
	return new(httpWriter)
}

type httpWriter struct {
}

func must(i string, err error) string {
	if err != nil {
		panic(err)
	}
	return i
}
func (h httpWriter) Write(p []byte) (n int, err error) {
	str := config.Get("dbg.httpUrl")
	if str == "" {
		return 0, nil
	}
	payload := fmt.Sprintf(string(p[:]))
	k := fmt.Sprintf("%v:%v:%v ", must(os.Hostname()), os.Getpid(), time.Now().UnixNano())
	str = fmt.Sprintf("http://%s?%s=%s", str, url.QueryEscape(k), url.QueryEscape(payload))
	resp, err := http.Get(str)
	if err != nil {
		return 0, fmt.Errorf("dbg.httpWriter:", err)
	}
	defer resp.Body.Close()
	return 0, nil
}

func makeFileWriter() io.Writer {
	return new(fileWriter)
}

type fileWriter struct {
}

func (h fileWriter) Write(p []byte) (n int, err error) {
	str := config.Get("dbg.logfile")
	if str == "" {
		return 0, nil
	}
	f, err := os.OpenFile(str, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("dbg.filewriter:", err)
	}
	defer f.Close()
	return f.Write(p)
}
