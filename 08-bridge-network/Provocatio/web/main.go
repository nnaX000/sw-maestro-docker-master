package main

import (
	"fmt"
	"io"
	"net/http"
)

// 사용자 요청을 받으면 api 서비스를 '이름'으로 호출해 다음 숫자를 가져온다.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://api:6000/next")
		if err != nil {
			http.Error(w, "api 호출 실패: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(w, "Hello from web! (visit #%s)\n", string(body))
	})
	http.ListenAndServe(":5000", nil)
}
