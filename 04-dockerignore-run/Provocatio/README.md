# 문제 04. .dockerignore와 RUN 통합 — 이미지 군살 빼기

**내용 범위:** 빌드 컨텍스트, `.dockerignore`, `RUN` 레이어 누적과 통합, `docker history`

> 빌드는 되지만 이미지가 필요 이상으로 무겁고, 비밀값(.env)이나 로그 같은 게 컨텍스트로 다 딸려 들어갑니다. 두 가지로 정리합니다.

## 상황

```
Provocatio/app/
├── main.go, go.mod, go.sum
├── .env             ← 비밀값 (이미지에 절대 들어가면 안 됨)
├── tmp/big.bin      ← 빌드와 무관한 큰 파일 (컨텍스트만 키움)
├── debug.log        ← 로그 파일
└── NOTES.local.md   ← 로컬 메모
```

아래는 **현재의 비효율적인 Dockerfile**입니다. 이 내용으로 `app/Dockerfile`을 만든 뒤 고치세요. apt 명령이 여러 `RUN`으로 쪼개져 있고, 받은 패키지 목록을 지우지 않습니다.

```dockerfile
FROM golang:1.22
WORKDIR /app
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .
EXPOSE 5000
CMD ["./server"]
```

## 요구사항

**(1) `.dockerignore` 작성** — `app/`에 `.dockerignore`를 만들어 `.env`, `tmp/`, `*.log`, `NOTES.local.md` 등 빌드에 불필요하거나 민감한 것을 컨텍스트에서 제외하세요.

**(2) RUN 통합** — 쪼개진 `apt-get update` / `install` 을 **하나의 `RUN`** 으로 합치고, 끝에 패키지 목록 캐시(`/var/lib/apt/lists/*`)를 정리하세요.

이미지 이름은 `counter:4.0`.

## 검증

```bash
# 1) 빌드 (컨텍스트 전송량이 줄었는지 첫 줄의 build context 크기를 보세요)
docker build -t counter:4.0 ./app

# 2) 레이어별 용량 확인 — apt 관련 레이어가 하나로 합쳐졌는지
docker history counter:4.0

# 3) 민감/불필요 파일이 이미지 안에 안 들어갔는지 확인
docker run --rm counter:4.0 ls -a /app
# → .env, tmp, debug.log, NOTES.local.md 가 보이지 않아야 함
```
