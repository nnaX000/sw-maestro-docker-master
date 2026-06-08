package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var (
		mu    sync.Mutex
		count int
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		c := count
		mu.Unlock()
		fmt.Fprintf(w, "Hello from Container! (visit #%d)\n", c)
	})
	http.ListenAndServe(":5000", nil)
}
