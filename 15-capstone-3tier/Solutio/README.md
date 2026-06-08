# 해설 — 문제 15 (캡스톤)

## 1. api/Dockerfile (multi-stage, 5번 적용)

```dockerfile
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM gcr.io/distroless/static
COPY --from=builder /app/server /server
EXPOSE 6000
ENTRYPOINT ["/server"]
```

## 2. web/Dockerfile (nginx 정적 + 프록시)

```dockerfile
FROM nginx:alpine
COPY index.html /usr/share/nginx/html/index.html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

`nginx.conf`는 `/`로 React 페이지를 서빙하고, `/api/`를 `http://api:6000/`으로 프록시합니다. 브라우저는 같은 출처(web)로만 요청하고, 실제 api 호출은 컨테이너 네트워크 안에서 **이름(`api`)으로** 일어납니다(8번).

## 3. 조립 순서와 명령

```bash
docker network create appnet
docker volume create pgdata

# (1) DB — 데이터는 named volume, 외부 노출 없음
docker run -d --name db --network appnet \
  -e POSTGRES_USER=app -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=app \
  -v pgdata:/var/lib/postgresql/data postgres:16

# (2) API — DB에 이름으로 연결(db:5432), 외부 노출 없음
docker build -t capstone-api:1.0 ./api
docker run -d --name api --network appnet \
  -e DATABASE_URL="postgres://app:secret@db:5432/app?sslmode=disable" \
  capstone-api:1.0

# (3) WEB — 유일하게 외부 노출(8080)
docker build -t capstone-web:1.0 ./web
docker run -d --name web --network appnet -p 8080:80 capstone-web:1.0
```

> api는 시작 시 DB가 아직 안 떠 있을 수 있어 `db.Ping()`을 몇 초간 재시도하도록 작성돼 있습니다. 그래서 db보다 먼저 떠도 곧 연결됩니다.

## 4. 각 tier가 어느 문제에서 왔나

| tier | 적용한 개념 | 출처 문제 |
|------|-------------|-----------|
| web (nginx) | 정적 서빙 + 외부 노출만 `-p` | 9번 |
| api (Go) | multi-stage 경량 이미지, ENTRYPOINT | 2·5번 |
| 통신 | user-defined network, 이름 기반(`db`, `api`) | 8번 |
| db (Postgres) | named volume 영속, 컨테이너 교체에도 데이터 유지 | 10번 |
| 보안 | 내부 서비스(api·db)는 미노출 | 9번 |

## 5. 영속성 검증

```bash
# 방문 수를 몇 번 올린 뒤, db 컨테이너만 통째로 교체
docker rm -f db
docker run -d --name db --network appnet \
  -e POSTGRES_USER=app -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=app \
  -v pgdata:/var/lib/postgresql/data postgres:16
# 다시 '방문하기' → 이전 숫자 다음부터 이어진다 (pgdata 볼륨 덕분)
```

1번에서 컨테이너를 지우면 0으로 돌아가던 그 카운터가, 이제 DB와 볼륨에 실려 **컨테이너가 죽고 다시 떠도 살아남습니다**. 시리즈의 두 축 — "무거운 이미지 → 경량화"와 "휘발성 → 영속성" — 이 여기서 완성됩니다.

## 6. 더 나아가기

- web도 multi-stage로(예: `node`로 빌드 → `nginx`로 서빙) 바꿔보기
- api에 `/healthz` 추가하고 `--health-cmd`로 헬스체크 붙이기
- 14번의 리소스 한도를 각 컨테이너에 적용해보기
