# 베이스는 golang:1.22, 작업 디렉터리는 /app
# 빌드는 1번과 동일하게 (go build -o server .)
# ENTRYPOINT 로 ./server를 항상 실행되게 고정한다 (exec form 사용)
# CMD 로 기본 인자 -msg "Hello from Container!" 를 제공한다 (exec form)
# 포트 5000 명시, 이미지 이름 counter:2.0

FROM golang:1.22
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o server .
EXPOSE 5000

ENTRYPOINT ["./server"]
CMD ["-msg", "Hello from Container"]

docker rm -f my-container-2
docker run -d --name my-container-2 -p 8081:5000 counter:2.0 -msg "Hello (prod)"
curl http://localhost:8081