# 해설 — 문제 08

## 1. Dockerfile (api / web 공통 형식)

각 디렉터리에 동일한 형식의 Dockerfile을 둡니다. (api는 6000, web은 5000 포트)

```dockerfile
# api/Dockerfile
FROM golang:1.22
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o server .
EXPOSE 6000
CMD ["./server"]
```

```dockerfile
# web/Dockerfile  (EXPOSE만 5000으로)
FROM golang:1.22
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o server .
EXPOSE 5000
CMD ["./server"]
```

```bash
docker build -t api:8.0 ./api
docker build -t web:8.0 ./web
```

## 2. 실행

```bash
docker network create appnet
docker run -d --name api --network appnet api:8.0
docker run -d --name web --network appnet -p 8080:5000 web:8.0
curl http://localhost:8080   # Hello from web! (visit #1)
```

## 3. 왜 이름으로 통신되는가 — embedded DNS

- **user-defined 네트워크**(여기선 `appnet`)에는 Docker 내장 DNS가 `127.0.0.11`에 떠 있습니다.
- 같은 네트워크에 붙은 컨테이너는 **컨테이너 이름 = 호스트명**으로 서로를 찾습니다. `web`이 `api`를 질의하면 DNS가 `api` 컨테이너의 현재 IP로 응답합니다.
- 그래서 `web` 코드의 `http://api:6000/next` 가 동작합니다. **컨테이너가 재시작되어 IP가 바뀌어도 이름은 그대로**라 안정적입니다.

## 4. 기본 bridge와의 차이

`docker0` 기본 bridge에는 이 이름 해석(DNS)이 **없습니다**. 그래서 기본 bridge에서 `web`을 띄우면 `api`라는 이름을 못 찾아 `no such host`로 실패합니다.

> 결론: 멀티 컨테이너 앱은 **반드시 user-defined 네트워크**를 만들어 이름으로 묶으세요. 이게 15번 3-tier의 기본 골격이 됩니다.

## 5. api는 왜 -p 가 없나

`api`는 외부(호스트)에서 직접 접근할 필요가 없고, 오직 `web`만 내부 네트워크로 호출합니다. 그래서 `-p`로 노출하지 않습니다. 외부로 열 필요가 없는 건 열지 않는 것이 보안의 기본입니다. (다음 9번 주제)
