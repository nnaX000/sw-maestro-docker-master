# 문제 08. 컨테이너 간 통신 — 이름이 곧 주소

**내용 범위:** user-defined bridge 네트워크, 컨테이너 이름 기반 통신(embedded DNS), 기본 bridge와의 차이

> 카운터를 두 조각으로 나눴습니다. 숫자를 보관하는 `api`와, 사용자를 받아 api를 호출하는 `web`. 둘은 서로의 IP를 모릅니다. **이름으로** 통신하게 만드세요. (15번 3-tier의 예고편)

## 상황

```
web (5000) ──HTTP──> api (6000)
```

`web`은 코드에서 `http://api:6000/next` 를 호출합니다. 즉 **`api`라는 이름으로** 접근합니다. IP가 아니라 이름으로 찾을 수 있어야 동작합니다.

## 주어지는 것

```
Provocatio/
├── api/   (main.go, go.mod)   ← /next 호출 시 다음 숫자 반환, :6000
└── web/   (main.go, go.mod)   ← http://api:6000/next 호출, :5000
```

## 요구사항

1. `api`, `web` 각각을 이미지로 빌드한다 (`api:8.0`, `web:8.0`) — Dockerfile은 1번과 같은 형식
2. **user-defined 네트워크**를 하나 만든다 (예: `appnet`)
3. 두 컨테이너를 **같은 네트워크에 이름과 함께** 띄운다 (api는 반드시 `--name api`)
4. `web`만 호스트 8080으로 노출한다 (api는 외부 노출 불필요)
5. (비교) 기본 bridge에서는 이름 통신이 안 되는 것도 확인한다

## 검증

```bash
docker network create appnet
docker run -d --name api --network appnet api:8.0
docker run -d --name web --network appnet -p 8080:5000 web:8.0

curl http://localhost:8080
# → Hello from web! (visit #1)
curl http://localhost:8080
# → Hello from web! (visit #2)

# (비교) 같은 네트워크가 아니면? — 이름 해석 실패
docker run --rm web:8.0 &   # 기본 bridge에서 실행하면 web이 api를 못 찾음
# web 응답: api 호출 실패 ... no such host
```
