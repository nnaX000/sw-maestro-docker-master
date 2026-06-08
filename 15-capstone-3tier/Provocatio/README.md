# 문제 15. 캡스톤 — 3-tier 스택 손으로 조립하기

**내용 범위:** 1~14번 종합 — 이미지 빌드·멀티스테이지·네트워크·이름 기반 통신·named volume·DB 영속성. **docker-compose 없이** plain CLI로만.

> 마지막입니다. 지금까지 배운 모든 것을 모아, 프론트엔드 + API + DB로 이뤄진 3-tier 카운터를 직접 조립합니다.

## 구성

```
브라우저 ──> web (nginx, React 정적 + /api 프록시) ──> api (Go) ──> db (Postgres + named volume)
```

- **web**: nginx가 React 정적 페이지를 서빙하고, `/api/` 요청을 `api` 컨테이너로 프록시
- **api**: Go 서버. `/next` 호출 시 Postgres의 카운터를 1 증가시켜 반환
- **db**: 공식 `postgres:16`, 데이터는 named volume에 영속

## 주어지는 것

```
Provocatio/
├── api/   (main.go, go.mod, go.sum)   ← lib/pq로 Postgres 연결, :6000, /next
└── web/   (index.html, nginx.conf)    ← React 프론트 + 프록시 설정
```

## 요구사항

1. **api 이미지**: multi-stage로 빌드 (`capstone-api:1.0`) — 5번에서 배운 distroless 경량화 적용
2. **web 이미지**: `nginx:alpine` 기반으로 `index.html`과 `nginx.conf`를 넣어 빌드 (`capstone-web:1.0`)
3. **네트워크**: user-defined 네트워크 `appnet` 생성, 세 컨테이너를 모두 여기에 (8번)
4. **db**: `postgres:16`을 `--name db`로 띄우고, 데이터는 named volume `pgdata`에 (10번)
   - 환경변수: `POSTGRES_USER=app`, `POSTGRES_PASSWORD=secret`, `POSTGRES_DB=app`
5. **api**: `DATABASE_URL` 환경변수로 db에 연결 (외부 노출 불필요 — 9번)
6. **web**: 호스트 8080 → 컨테이너 80으로 노출
7. **검증**: 브라우저에서 방문 수를 올린 뒤, **db 컨테이너를 삭제·재생성해도 숫자가 유지**되는지 확인 (volume 영속성)

## 검증 (예시 흐름)

```bash
docker network create appnet
docker volume create pgdata

# DB
docker run -d --name db --network appnet \
  -e POSTGRES_USER=app -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=app \
  -v pgdata:/var/lib/postgresql/data postgres:16

# API
docker build -t capstone-api:1.0 ./api
docker run -d --name api --network appnet \
  -e DATABASE_URL="postgres://app:secret@db:5432/app?sslmode=disable" capstone-api:1.0

# WEB
docker build -t capstone-web:1.0 ./web
docker run -d --name web --network appnet -p 8080:80 capstone-web:1.0

# 브라우저로 http://localhost:8080 → '방문하기' 클릭 → 숫자 증가
# DB 영속성 검증: db만 교체해도 숫자 유지
docker rm -f db
docker run -d --name db --network appnet \
  -e POSTGRES_USER=app -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=app \
  -v pgdata:/var/lib/postgresql/data postgres:16
# 다시 방문 → 이전 숫자 다음부터 이어짐
```
