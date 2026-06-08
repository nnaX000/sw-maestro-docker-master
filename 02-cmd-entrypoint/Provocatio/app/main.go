package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	msg := flag.String("msg", "Hello", "인사 메시지")
	flag.Parse()

	var (
		mu    sync.Mutex
		count int
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		c := count
		mu.Unlock()
		fmt.Fprintf(w, "%s (visit #%d)\n", *msg, c)
	})
	http.ListenAndServe(":5000", nil)
}
