package trickle

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type trickleReader struct {
	src  io.Reader
	ctx  context.Context
	rate Rate
}

func (t *trickleReader) Read(p []byte) (n int, err error) {
	if t.rate.Bytes < len(p) {
		p = p[0:t.rate.Bytes]
	}
	select {
	case <-t.ctx.Done():
		err = io.EOF
	case <-time.After(t.rate.Interval):
		n, err = t.src.Read(p)
	}
	return
}

// Rate describes a rate at which to read; X bytes after Y time.
type Rate struct {
	Bytes    int
	Interval time.Duration
}

// Reader creates a new io.Reader that will read at the specified rate, returning an
// io.EOF immediately on the given ctx being cancelled or timed out.
func Reader(src io.Reader, ctx context.Context, rate Rate) (io.Reader, error) {
	if rate.Bytes < 1 {
		return nil, fmt.Errorf("bytes must be > 0")
	}
	if rate.Interval < 1 {
		return nil, fmt.Errorf("interval must be > 0")
	}
	return &trickleReader{src, ctx, rate}, nil
}

type fileStreamer struct {
	data []byte
	rate Rate
}

// FileStreamer reads a file once and then streams its contents to each new http request.
func FileStreamer(path string, rate Rate) (http.Handler, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file data: %s", err)
	}
	return &fileStreamer{data, rate}, nil
}

// ServeHTTP implements the http.Handler interface.
func (rr *fileStreamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reader, err := Reader(bytes.NewReader(rr.data), r.Context(), rr.rate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, reader)
}
