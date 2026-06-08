package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	dataFile   = "/data/count"
	configFile = "/config/message.txt"
)

func load() int {
	b, err := os.ReadFile(dataFile)
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(b)))
	return n
}

func greeting() string {
	b, err := os.ReadFile(configFile)
	if err != nil {
		return "Hello (default)"
	}
	return strings.TrimSpace(string(b))
}

func main() {
	var mu sync.Mutex
	count := load()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		c := count
		_ = os.WriteFile(dataFile, []byte(strconv.Itoa(c)), 0644)
		mu.Unlock()
		// 설정 파일은 매 요청마다 읽으므로, 호스트에서 고치면 즉시 반영된다
		fmt.Fprintf(w, "%s (visit #%d)\n", greeting(), c)
	})
	http.ListenAndServe(":5000", nil)
}
