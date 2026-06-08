# 해설 — 문제 11

## 1. 실행

```dockerfile
# Dockerfile — /data 와 /config 를 미리 생성
FROM golang:1.22
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o server .
RUN mkdir -p /data /config
EXPOSE 5000
CMD ["./server"]
```

```bash
mkdir -p ~/myconfig && echo "Hello (bind!)" > ~/myconfig/message.txt
docker run -d --name {{my-container}} -p 8080:5000 \
  -v countdata:/data \
  -v ~/myconfig/message.txt:/config/message.txt \
  counter:11.0
```

## 2. Volume vs Bind Mount

| | Volume (`-v countdata:/data`) | Bind Mount (`-v ~/myconfig/...:/config/...`) |
|---|---|---|
| 관리 주체 | Docker (`/var/lib/docker/volumes`) | 사용자 (호스트의 특정 경로) |
| 경로 | 신경 안 써도 됨 | 호스트 절대경로를 직접 지정 |
| 이식성 | 좋음 (백업·이관 쉬움) | 호스트 구조에 의존 |
| 주 용도 | **운영 데이터·DB 영속** | **개발·설정 주입·로컬 소스** |

- bind mount는 호스트 경로를 컨테이너 경로에 **그대로 연결**합니다. 그래서 호스트에서 파일을 고치면 컨테이너가 보는 내용도 즉시 바뀝니다(앱이 매 요청마다 읽으므로 재시작도 불필요).
- 반대로 named volume은 "어디에 저장되는지" 신경 쓰지 않고 Docker에 맡기는 방식이라 운영 데이터에 적합합니다.

## 3. 언제 무엇을

- **설정 파일 주입 / 개발 중 소스 즉시 반영** → bind mount
- **운영 데이터 / DB / 영속 저장** → named volume

이 문제처럼 둘을 함께 쓰는 게 일반적입니다: 데이터는 볼륨, 설정은 bind mount.
