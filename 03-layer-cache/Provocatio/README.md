# 문제 03. 레이어 캐시 최적화 — COPY 순서 한 줄의 차이

**내용 범위:** 이미지 레이어와 빌드 캐시, `COPY` 순서, 캐시 무효화, `docker build` 캐시 관찰

> 카운터에 의존성(`github.com/google/uuid`)이 하나 생겼습니다. 이제 빌드 때마다 의존성을 새로 받느라 느려질 수 있습니다. Dockerfile 한 줄의 순서로 이걸 해결합니다.

## 상황

코드 한 줄만 고쳐도 빌드가 매번 처음부터 다 돕니다.

> "main.go에서 메시지 문구 하나 바꿨을 뿐인데 의존성 다운로드부터 전부 다시 도네요. 왜 이렇게 느리죠?"

## 주어지는 것

```
Provocatio/
└── app/
    ├── main.go
    └── go.mod / go.sum   ← uuid 의존성
```

아래는 **현재의 비효율적인 Dockerfile**입니다. 이 내용으로 `app/Dockerfile`을 만든 뒤, 이걸 고치는 것이 과제입니다. 소스 전체를 먼저 복사한 탓에 코드 한 줄만 바꿔도 의존성을 다시 받습니다.

```dockerfile
FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server .
EXPOSE 5000
CMD ["./server"]
```

## 요구사항

`Dockerfile`을 수정해서, **소스 코드만 바뀌었을 때 의존성 다운로드 단계가 캐시로 재사용**되게 만드세요.

1. `go.mod`, `go.sum`을 **먼저** 복사하고 `go mod download` 실행
2. 그 **다음에** 나머지 소스를 복사하고 `go build`
3. 이미지 이름 `counter:3.0`

## 검증

```bash
# 1) 최초 빌드
docker build -t counter:3.0 ./app

# 2) main.go의 메시지 문구를 아무거나 살짝 바꾼다 (예: Hello → Hi)
#    그리고 다시 빌드
docker build -t counter:3.0 ./app

# 3) 두 번째 빌드 로그를 본다
#  → COPY go.mod go.sum / RUN go mod download 단계에
#    CACHED 가 찍히면 성공 (의존성을 다시 안 받음)
#  → 비효율 버전이었다면 이 단계부터 전부 다시 실행된다
```

> 의존성이 수십 개로 늘어난 실제 프로젝트에서는 이 순서 한 줄이 빌드 시간을 몇 배 가릅니다.
