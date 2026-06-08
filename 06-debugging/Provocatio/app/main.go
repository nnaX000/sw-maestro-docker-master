package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	greeting := os.Getenv("GREETING")
	if greeting == "" {
		// 환경변수가 없으면 시작하자마자 종료된다
		log.Fatal("FATAL: GREETING 환경변수가 비어 있습니다. -e GREETING=... 로 전달하세요")
	}

	var (
		mu    sync.Mutex
		count int
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		c := count
		mu.Unlock()
		fmt.Fprintf(w, "%s (visit #%d)\n", greeting, c)
	})
	log.Println("listening on :5000")
	http.ListenAndServe(":5000", nil)
}
