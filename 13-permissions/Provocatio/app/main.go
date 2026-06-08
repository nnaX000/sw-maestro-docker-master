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

func main() {
	var mu sync.Mutex
	count := load()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		c := count
		err := os.WriteFile(dataFile, []byte(strconv.Itoa(c)), 0644)
		mu.Unlock()
		if err != nil {
			// 권한이 없으면 여기서 에러가 드러난다
			http.Error(w, "쓰기 실패: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "visit #%d\n", c)
	})
	http.ListenAndServe(":5000", nil)
}
