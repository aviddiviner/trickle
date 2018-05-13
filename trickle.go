// Package trickle has some simple functions for streaming data at a chosen rate.
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

func (r Rate) valid() (err error) {
	if r.Bytes < 1 {
		err = fmt.Errorf("bytes must be > 0")
	}
	if r.Interval < 1 {
		err = fmt.Errorf("interval must be > 0")
	}
	return
}

// Reader wraps a source, creating a new io.Reader that will read at a specified rate.
// Reader returns io.EOF immediately on the given ctx being cancelled or timing out.
func Reader(source io.Reader, ctx context.Context, rate Rate) (io.Reader, error) {
	if err := rate.valid(); err != nil {
		return nil, fmt.Errorf("invalid rate: %s", err)
	}
	return &trickleReader{source, ctx, rate}, nil
}

type fileStreamer struct {
	data []byte
	rate Rate
}

// FileStreamer reads a file once and then streams its contents on each new http request.
func FileStreamer(path string, rate Rate) (http.Handler, error) {
	if err := rate.valid(); err != nil {
		return nil, fmt.Errorf("invalid rate: %s", err)
	}
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
