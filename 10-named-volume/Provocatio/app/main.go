package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const dataFile = "/data/count"

func load() int {
	b, err := os.ReadFile(dataFile)
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(b)))
	return n
}

func save(n int) error {
	return os.WriteFile(dataFile, []byte(strconv.Itoa(n)), 0644)
}

func main() {
	var mu sync.Mutex
	count := load() // 시작할 때 디스크에서 복원

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		c := count
		_ = save(c)
		mu.Unlock()
		fmt.Fprintf(w, "visit #%d\n", c)
	})
	http.ListenAndServe(":5000", nil)
}
