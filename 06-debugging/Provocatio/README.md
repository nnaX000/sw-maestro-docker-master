# 문제 06. 컨테이너 디버깅 — 안 뜨는 컨테이너 추적하기

**내용 범위:** `docker logs`, `docker ps -a`, `docker inspect`, `docker exec`, `docker cp` 로 장애 진단

> 동료가 카운터를 배포했는데 컨테이너가 자꾸 죽는다고 합니다. 코드는 못 고친다고 가정하고, 진단 명령만으로 원인을 찾아 정상 기동시키세요.

## 상황

이미지를 빌드하고 평소처럼 띄웠는데 곧바로 죽습니다.

```bash
docker build -t counter:6.0 ./app
docker run -d --name {{my-container}} -p 8080:5000 counter:6.0
curl http://localhost:8080
# → Connection refused ... 분명 run 은 성공했는데?
```

`main.go`는 수정하지 않습니다. **왜 죽는지 진단하고**, 컨테이너가 계속 떠 있도록 **올바른 실행 명령**을 찾는 것이 과제입니다.

## 주어지는 것

```
Provocatio/app/   (main.go, go.mod)
```

`app/`에 아래 내용으로 `Dockerfile`을 만들어 빌드하세요. (Dockerfile 자체는 정상입니다 — 버그는 다른 데 있습니다.)

```dockerfile
FROM golang:1.22
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o server .
EXPOSE 5000
CMD ["./server"]
```

## 요구사항

아래 진단 명령들을 활용해 원인을 찾고, 카운터를 정상 기동시키세요.

1. `docker ps -a` 로 컨테이너 상태(STATUS)를 확인
2. `docker logs {{my-container}}` 로 종료 원인 메시지를 확인
3. `docker inspect` 로 종료 코드(ExitCode) 확인
4. 원인을 해결한 올바른 `docker run` 명령으로 다시 기동
5. 기동 후 `docker exec -it`로 들어가 프로세스가 살아 있는지 확인

## 검증

```bash
# 정상 기동 후
curl http://localhost:8080
# → <당신이 정한 인사말> (visit #1)

docker ps
# → STATUS가 Up (Exited 가 아님)
```
