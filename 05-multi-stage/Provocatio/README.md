# 문제 05. Multi-stage build — 이미지를 1/50로

**내용 범위:** multi-stage build, 빌드 단계와 런타임 단계 분리, 정적 바이너리, distroless/scratch

> 1번부터 들고 온 그 무거운 이미지(수백 MB)를 기억하시죠? 이번에 10MB대로 줄입니다. Go가 단일 바이너리로 컴파일된다는 성질을 활용합니다.

## 상황

`docker images counter` 를 보면 지금까지의 이미지는 수백 MB입니다. 안에는 Go 컴파일러, 소스, 빌드 캐시가 전부 들어 있습니다. 하지만 **실행에 필요한 건 컴파일된 바이너리 하나**뿐입니다.

## 요구사항

`app/`에 multi-stage `Dockerfile`을 작성하세요.

1. **빌드 단계**: `golang:1.22 AS builder` 에서 정적 바이너리로 컴파일
   - 정적 링크를 위해 `CGO_ENABLED=0` 으로 빌드 (`RUN CGO_ENABLED=0 go build -o server .`)
2. **런타임 단계**: `gcr.io/distroless/static` (또는 `scratch`) 를 베이스로
   - 빌드 단계의 `server` 바이너리만 `COPY --from=builder` 로 가져온다
3. 포트 5000 명시, 진입점은 바이너리, 이미지 이름 `counter:5.0`

## 검증

```bash
docker build -t counter:5.0 ./app

# 용량 비교 — 극적으로 줄었는지
docker images | grep counter
# → counter  5.0  ...  약 10~15MB   (이전 버전은 수백 MB)

# 그래도 똑같이 동작하는지
docker run -d --name {{my-container}} -p 8080:5000 counter:5.0
curl http://localhost:8080
# → Hello (visit #1, req xxxxxxxx)
```
