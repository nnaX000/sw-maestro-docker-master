# 해설 — 문제 04

## 1. .dockerignore 정답

`app/.dockerignore`:

```
# 민감/불필요 파일은 빌드 컨텍스트에서 제외
.env
tmp/
*.log
*.local.md
.git
Dockerfile
.dockerignore
```

## 2. Dockerfile 정답

```dockerfile
FROM golang:1.22
WORKDIR /app

# 여러 RUN을 하나로 합치고, 같은 레이어에서 정리까지 끝낸다
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

EXPOSE 5000
CMD ["./server"]
```

## 3. .dockerignore — 왜 필요한가

`docker build ./app` 은 `app/` 디렉터리 **전체**를 데몬에 빌드 컨텍스트로 전송합니다. `.git`이나 로그가 크면 전송이 느려지고, `.env` 같은 민감 파일이 이미지에 섞여 들어갈 위험도 있습니다. `.dockerignore`로 처음부터 제외하면 전송도 빨라지고 사고도 막습니다.

## 4. RUN 통합 — 왜 한 줄로

이미지는 레이어의 **누적**입니다. 한 번 쌓인 레이어는 불변이라, 다음 `RUN`에서 파일을 지워도 **이전 레이어엔 그대로 남아** 용량에 포함됩니다.

```dockerfile
# 나쁨 — 3개 레이어, apt 캐시가 이미지에 박힘
RUN apt-get update
RUN apt-get install -y ca-certificates
# (목록 캐시가 어딘가 레이어에 남음)

# 좋음 — 1개 레이어, 받자마자 같은 레이어에서 정리
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
```

`docker history`로 레이어별 용량을 보면 차이가 드러납니다. "최종 파일시스템에 안 보여도 거쳐온 레이어 합이 곧 용량"이라는 점이 핵심입니다.
