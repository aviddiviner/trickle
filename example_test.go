package trickle_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aviddiviner/trickle"
)

func ExampleReader() {
	// ---- Create a trickle.Reader ----

	var src = bytes.NewBufferString("Hello world") // 11 bytes
	var ctx = context.TODO()

	r, _ := trickle.Reader(src, ctx, trickle.Rate{
		Bytes:    4,               // 4 bytes ...
		Interval: 1 * time.Second, // per second
	})

	// ---- Test reading from it ----

	var buf = make([]byte, 11)

	readBytes := func(r io.Reader, buf []byte) {
		if n, _ := r.Read(buf); n > 0 {
			fmt.Printf("read %d bytes: %q\n", n, buf[0:n])
		}
	}

	readBytes(r, buf) // takes 1 second
	readBytes(r, buf) // takes 1 second
	readBytes(r, buf) // takes 1 second

	// Output:
	// read 4 bytes: "Hell"
	// read 4 bytes: "o wo"
	// read 3 bytes: "rld"
}
