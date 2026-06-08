# 해설 — 문제 03

## 1. Dockerfile 정답

```dockerfile
FROM golang:1.22
WORKDIR /app

# 의존성 정의 먼저 복사 → 다운로드  (이 레이어가 캐시의 핵심)
COPY go.mod go.sum ./
RUN go mod download

# 그다음 소스 복사 → 빌드
COPY . .
RUN go build -o server .

EXPOSE 5000
CMD ["./server"]
```

## 2. 왜 순서가 중요한가

Docker는 Dockerfile의 **각 명령을 레이어 하나**로 만들고, 레이어 단위로 캐시합니다. 어떤 레이어의 입력이 바뀌면 **그 레이어와 이후 모든 레이어**의 캐시가 무효화됩니다.

- **나쁜 순서** (`COPY . .` 먼저): `main.go`를 한 글자만 고쳐도 `COPY . .` 레이어의 입력이 바뀝니다 → 그 아래 `go mod download`, `go build`가 **전부 다시** 실행됩니다.
- **좋은 순서** (`COPY go.mod go.sum` 먼저): 소스만 바뀌면 `go.mod`/`go.sum`은 그대로이므로 `go mod download` 레이어는 **CACHED**. 변경된 소스 때문에 다시 도는 건 `COPY . .`와 `go build`뿐입니다.

## 3. 원칙

> **잘 안 변하는 것을 위로, 자주 변하는 것을 아래로.**

의존성 정의(go.mod/go.sum, package.json, requirements.txt)는 드물게 바뀌고 소스는 자주 바뀝니다. 그래서 의존성 설치를 위쪽 레이어로 분리하는 것이 모든 언어의 공통 패턴입니다.

## 4. 캐시 확인 팁

- 두 번째 빌드 로그에서 `CACHED [.. RUN go mod download]` 가 보이면 성공.
- `docker build --no-cache` 로 캐시를 무시하고 처음부터 빌드해 비교해볼 수 있습니다.
