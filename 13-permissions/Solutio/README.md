# 해설 — 문제 13

## 1. Dockerfile 정답 (비root 실행)

```dockerfile
FROM golang:1.22
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o server .
RUN mkdir -p /data && \
    useradd -u 1000 -m appuser && \
    chown -R 1000:1000 /data
USER 1000
EXPOSE 5000
CMD ["./server"]
```

이미지 내부의 `/data`는 uid 1000 소유로 만들어 둡니다. 문제는 **bind mount 시 호스트 디렉터리의 소유자가 우선**한다는 점입니다.

## 2. 왜 permission denied 인가

- 컨테이너 안 프로세스는 uid 1000으로 동작합니다.
- bind mount는 호스트 디렉터리를 그대로 연결하는데, 그 디렉터리가 **root(uid 0)** 또는 다른 uid 소유라면, uid 1000 프로세스는 쓰기 권한이 없습니다 → `permission denied`.
- 즉 "컨테이너 프로세스 UID ≠ 호스트 디렉터리 소유자 UID" 가 단골 원인입니다.

## 3. 해결 방법 세 가지

| 방법 | 명령 | 비고 |
|------|------|------|
| 호스트 소유자 변경 | `sudo chown -R 1000:1000 ~/hostdata` | 가장 직접적 |
| 실행 UID 맞추기 | `docker run --user $(id -u):$(id -g) ...` | 호스트 디렉터리 소유자에 맞춤 |
| named volume 사용 | `-v countdata2:/data` | Docker가 권한을 비교적 자동 처리 |

## 4. 실무 권장

- bind mount에서 특히 권한 문제가 잦습니다. **운영 데이터는 named volume**을 쓰면 이 부담이 크게 줍니다.
- 이미지가 어떤 UID로 도는지 확인하세요(예: postgres 공식 이미지는 uid 999). 호스트 디렉터리를 그 UID로 맞추거나 named volume으로 회피합니다.
