# 문제 09. 포트 매핑과 바인딩 보안 — -p 를 제대로 이해하기

**내용 범위:** `-p`의 동작(DNAT), 바인딩 주소(0.0.0.0 vs 127.0.0.1), `EXPOSE`에 대한 오해

> `-p`를 둘러싼 흔한 착각 세 가지를 직접 실험으로 깨봅니다. 카운터 단일 컨테이너로 진행합니다.

## 주어지는 것

```
Provocatio/app/   (main.go, go.mod)   ← :5000 단순 카운터
```

먼저 `counter:9.0` 으로 빌드하세요. (Dockerfile은 1번 형식)

## 요구사항 — 아래 세 가지를 실험으로 확인하고 해설과 대조

**(1) EXPOSE만으로는 안 열린다**
Dockerfile에 `EXPOSE 5000`이 있어도, `-p` 없이 띄우면 호스트에서 접속되지 않음을 확인.

**(2) -p 의 방향과 동작**
`-p 8080:5000` 으로 띄우면 호스트 8080 → 컨테이너 5000 으로 연결됨을 확인.

**(3) 바인딩 주소 차이 (보안)**
`-p 8080:5000` (= `0.0.0.0:8080`, 외부 전부 노출) vs `-p 127.0.0.1:8080:5000` (로컬만) 의 차이를 확인.

## 검증

```bash
docker build -t counter:9.0 ./app

# (1) -p 없이
docker run -d --name c1 counter:9.0
curl http://localhost:8080            # → 연결 거부 (EXPOSE만으론 안 열림)

# (2) -p 로 노출
docker rm -f c1
docker run -d --name c1 -p 8080:5000 counter:9.0
curl http://localhost:8080            # → Hello ...
docker port c1                         # 5000/tcp -> 0.0.0.0:8080

# (3) 로컬 전용 바인딩
docker rm -f c1
docker run -d --name c1 -p 127.0.0.1:8080:5000 counter:9.0
curl http://localhost:8080            # → 로컬에선 OK
docker port c1                         # 5000/tcp -> 127.0.0.1:8080  (외부 차단)
```
