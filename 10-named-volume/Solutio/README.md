# 해설 — 문제 10

## 1. Dockerfile 정답

```dockerfile
FROM golang:1.22
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o server .
RUN mkdir -p /data          # 볼륨 없이 떠도 경로는 존재하게
EXPOSE 5000
CMD ["./server"]
```

## 2. 실행과 검증

```bash
docker volume create countdata
docker run -d --name {{my-container}} -p 8080:5000 -v countdata:/data counter:10.0
curl http://localhost:8080   # visit #1 ... #3

docker rm -f {{my-container}}
docker run -d --name {{my-container}} -p 8080:5000 -v countdata:/data counter:10.0
curl http://localhost:8080   # visit #4  → 데이터 유지
```

## 3. 왜 사라지고, 왜 살아남는가

- 컨테이너 = 이미지(읽기전용) + 그 위 얇은 **writable layer** 1장. 모든 변경은 writable layer에 기록됩니다.
- `docker rm` 은 이 writable layer를 삭제합니다 → 그 안에 있던 `/data/count`도 함께 소멸. 그래서 볼륨 없이 띄우면 재생성 시 0부터 시작합니다.
- **named volume**(`countdata`)은 컨테이너 **바깥**(`/var/lib/docker/volumes/...`)에 데이터를 둡니다. `-v countdata:/data` 로 마운트하면 `/data`에 쓰는 내용이 볼륨에 저장되어, 컨테이너를 지웠다 다시 만들어도 **같은 볼륨을 붙이는 한** 데이터가 유지됩니다.

## 4. 핵심 원칙

> "컨테이너는 언제든 죽고 다시 뜬다"를 전제로 설계하고, **영속이 필요한 데이터는 writable layer 바깥(볼륨)** 에 둔다.

DB 컨테이너를 운영할 때도 똑같습니다. 엔진별 데이터 경로(postgres=`/var/lib/postgresql/data`, mysql=`/var/lib/mysql`)에 볼륨을 정확히 매핑해야 컨테이너를 교체해도 데이터가 남습니다. (15번 캡스톤에서 실제로 적용)
