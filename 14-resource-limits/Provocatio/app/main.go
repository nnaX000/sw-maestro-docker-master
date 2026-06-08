package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
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
		fmt.Fprintf(w, "visit #%d\n", c)
	})

	// CPU 부하: ?sec=N 초 동안 모든 코어를 바쁘게 돌린다
	http.HandleFunc("/cpu", func(w http.ResponseWriter, r *http.Request) {
		sec, _ := strconv.Atoi(r.URL.Query().Get("sec"))
		if sec <= 0 {
			sec = 10
		}
		deadline := time.Now().Add(time.Duration(sec) * time.Second)
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for time.Now().Before(deadline) {
				}
			}()
		}
		fmt.Fprintf(w, "burning CPU for %ds on %d cores\n", sec, runtime.NumCPU())
	})

	// 메모리 부하: ?mb=N 만큼 실제로 할당해 붙잡는다 (한도 초과 시 OOM kill 유도)
	var hold [][]byte
	http.HandleFunc("/mem", func(w http.ResponseWriter, r *http.Request) {
		mb, _ := strconv.Atoi(r.URL.Query().Get("mb"))
		if mb <= 0 {
			mb = 100
		}
		for i := 0; i < mb; i++ {
			b := make([]byte, 1024*1024)
			for j := range b {
				b[j] = 1 // 실제로 메모리를 점유하도록 touch
			}
			hold = append(hold, b)
		}
		fmt.Fprintf(w, "allocated ~%dMB this call (held total ~%dMB)\n", mb, len(hold))
	})

	http.ListenAndServe(":5000", nil)
}
