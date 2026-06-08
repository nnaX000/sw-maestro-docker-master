# 해설 — 문제 07

## 1. 전체 흐름

```bash
# (1) 로컬 레지스트리 컨테이너 — 호스트 5001 → 레지스트리 5000
docker run -d -p 5001:5000 --name registry registry:2

# (2) 레지스트리 주소를 포함한 이름으로 태그
docker tag counter:5.0 localhost:5001/counter:5.0

# (3) push
docker push localhost:5001/counter:5.0

# (4) 로컬 캐시 삭제 후 pull 로 검증
docker rmi counter:5.0 localhost:5001/counter:5.0
docker pull localhost:5001/counter:5.0
docker run -d --name {{my-container}} -p 8080:5000 localhost:5001/counter:5.0
curl http://localhost:8080
```

## 2. 이미지 이름 규칙

이미지 전체 이름은 `[Registry]/[Repository]:[Tag]` 세 부분입니다.

| 부분 | 예 | 생략하면 |
|------|----|----------|
| Registry | `localhost:5001` | Docker Hub(`docker.io`)로 간주 |
| Repository | `counter` | (생략 불가, 단 공식은 `library/` 자동) |
| Tag | `5.0` | `latest` 자동 적용 |

- `docker pull ubuntu` 는 내부적으로 `docker.io/library/ubuntu:latest` 로 보완됩니다.
- **어디로 push할지는 이미지 이름의 Registry 부분이 결정**합니다. 그래서 push 전에 반드시 레지스트리 주소를 포함해 `docker tag` 해야 합니다.

## 3. Docker Hub로 올릴 때

```bash
docker login
docker tag counter:5.0 <내아이디>/counter:5.0
docker push <내아이디>/counter:5.0
```

개인 저장소 이름은 보통 `<사용자명>/<이미지>` 형태입니다. 공식(`library/...`) 저장소에는 개인이 push할 수 없습니다.
