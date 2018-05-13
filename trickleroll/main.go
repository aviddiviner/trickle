package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aviddiviner/trickle"
)

func main() {
	rickroller, err := trickle.FileStreamer("rickroll.mp4", trickle.Rate{
		Bytes:    2 << 10, // 2 << 10 == 2k
		Interval: 1 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	http.Handle("/", rickroller)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
