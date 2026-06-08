# 해설 — 문제 05

## 1. Dockerfile 정답

```dockerfile
# ---------- 1단계: 빌드 ----------
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# CGO를 끄면 외부 libc에 의존하지 않는 정적 바이너리가 나온다
RUN CGO_ENABLED=0 go build -o server .

# ---------- 2단계: 런타임 (최소) ----------
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 5000
ENTRYPOINT ["./server"]
```

> `scratch`(완전 빈 이미지)도 가능합니다. distroless/static은 비어 있되 CA 인증서·타임존 등 최소한이 들어 있어 더 실용적입니다.

## 2. 왜 이렇게 줄어드는가

- **빌드 단계**의 `golang:1.22`에는 컴파일러·표준 라이브러리 소스·빌드 도구가 전부 들어 있어 수백 MB입니다. 하지만 이건 **컴파일할 때만** 필요합니다.
- `COPY --from=builder` 는 빌드 단계의 결과물(`server` 바이너리) **하나만** 최종 이미지로 가져옵니다. 컴파일러는 빠집니다.
- 최종 이미지 = 빈 베이스 + 바이너리 한 개 → 10MB 안팎.

## 3. CGO_ENABLED=0 가 핵심

`distroless/static`이나 `scratch`에는 시스템 라이브러리(libc)가 없습니다. 기본 빌드는 동적 링크라 libc를 찾다가 실행 실패합니다. `CGO_ENABLED=0` 으로 **정적 바이너리**를 만들면 의존하는 외부 라이브러리가 없어 빈 이미지에서도 잘 뜹니다.

## 4. 일반화

multi-stage는 "빌드에만 필요한 것"과 "실행에 필요한 것"을 분리하는 패턴입니다. Go뿐 아니라 Node(빌드 후 dist만), Java(JDK로 빌드 → JRE로 실행) 등에도 똑같이 적용됩니다. Go는 산출물이 단일 정적 바이너리라 효과가 가장 극적입니다.
