# 문제 01. 첫 컨테이너 — 이미지 빌드하고 띄우기

**내용 범위:** Dockerfile 기본 명령어(FROM·WORKDIR·COPY·RUN·EXPOSE·CMD), `docker build`, 포트 매핑(`-p`), `docker exec -it`, 동작 검증

> 이 실습 시리즈는 하나의 앱(**방문 카운터**)을 1번부터 15번까지 한 단계씩 고도화합니다. 1번에서는 가장 먼저, 이 앱을 컨테이너로 띄우는 것부터 시작합니다.
>
> 문서의 `{{ }}`는 **여러분이 직접 정해서 채우는 값**입니다. 예: `{{my-container}}` → `web`

## 상황

여러분은 백엔드 팀에 막 합류했습니다. 동료가 작은 Go 웹 서버 코드(`main.go`)를 던져주며 말합니다.

> "이거 내 노트북에선 잘 도는데, 옆 사람 환경에선 Go 버전이 안 맞네 뭐네 하면서 또 안 돈대요. Docker 이미지로 만들어서 어디서든 똑같이 뜨게 해줄 수 있어요?"

여러분의 임무는 이 앱을 **컨테이너로 패키징해서 실행**하는 것입니다. 코드는 건드리지 않습니다. 오직 Dockerfile만 작성합니다.

## 주어지는 것

```
Provocatio/
├── README.md          ← 지금 보는 문서
└── app/
    ├── main.go        ← 수정 금지
    └── go.mod         ← 수정 금지
```

해설은 `Solutio/` 디렉터리에 있습니다. 먼저 직접 풀어본 뒤 열어보세요.

`main.go` (참고용, 수정하지 마세요)

```go
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
```

`go.mod`

```
module counter

go 1.22
```

> 접속할 때마다 방문 횟수가 1씩 늘어납니다. 이 "숫자"가 앞으로 중요한 역할을 합니다.

## 요구사항

`app/` 안에 `Dockerfile`을 직접 작성해서 아래를 만족시키세요.

1. 베이스 이미지는 `golang:1.22`를 사용한다
2. 작업 디렉터리는 `/app`으로 지정한다
3. `go.mod`를 먼저 복사하고, 그다음 소스를 복사한다
4. `go build`로 실행 바이너리를 만든다 (결과물 이름: `server`)
5. 컨테이너가 사용하는 포트(5000)를 명시한다
6. 컨테이너가 뜨면 자동으로 그 바이너리가 실행된다
7. 이미지 이름은 `counter:1.0`으로 빌드한다
8. 호스트 **8080** 포트로 접속하면 컨테이너 **5000** 포트로 연결된다

## 검증

직접 작성한 Dockerfile로 이미지를 빌드하고 컨테이너를 띄운 뒤, 아래가 모두 충족되면 성공입니다.

```bash
# 1) 이미지가 만들어졌는가
docker images counter
# → REPOSITORY  TAG   ...   counter  1.0

# 2) 컨테이너가 떠 있는가
docker ps
# → STATUS가 Up, PORTS에 0.0.0.0:8080->5000/tcp

# 3) 포트 매핑이 의도대로인가
docker port {{my-container}}
# → 5000/tcp -> 0.0.0.0:8080

# 4) 응답이 오는가 (여러 번 호출해보세요)
curl http://localhost:8080
# → Hello from Container! (visit #1)
curl http://localhost:8080
# → Hello from Container! (visit #2)

# 5) 컨테이너 안으로 들어가 빌드 결과물이 있는지 확인
docker exec -it {{my-container}} sh
/app # ls
# → server  main.go  go.mod   ← server 바이너리가 보이면 성공
/app # exit
```
