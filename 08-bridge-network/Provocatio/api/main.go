package main

import (
	"fmt"
	"net/http"
	"sync"
)

// 카운터를 보관하는 내부 서비스. 호출될 때마다 다음 숫자를 돌려준다.
func main() {
	var (
		mu    sync.Mutex
		count int
	)
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		c := count
		mu.Unlock()
		fmt.Fprintf(w, "%d", c)
	})
	http.ListenAndServe(":6000", nil)
}
