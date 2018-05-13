# trickle

[![GoDoc](https://godoc.org/github.com/aviddiviner/trickle?status.svg)](https://godoc.org/github.com/aviddiviner/trickle)

A simple Go package for streaming data at a limited rate.

See [trickleroll](https://github.com/aviddiviner/trickle/tree/master/trickleroll) for an example webserver which responds with [rick rolls](https://www.youtube.com/watch?v=dQw4w9WgXcQ) at a couple of bytes per second:

	rickroller, _ := trickle.FileStreamer("rickroll.mp4", trickle.Rate{
		Bytes:    2 << 10, // 2k
		Interval: 1 * time.Second,
	})
	http.Handle("/", rickroller)
	log.Fatal(http.ListenAndServe(":8080", nil))

See [godoc.org](https://godoc.org/github.com/aviddiviner/trickle) for the rest of the package documentation.

## License

[MIT](LICENSE)
