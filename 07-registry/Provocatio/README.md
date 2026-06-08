# 문제 07. 레지스트리 push/pull — 이미지 공유하기

**내용 범위:** 이미지 naming 규칙(`[registry]/[repo]:[tag]`), `docker tag`, `push`, `pull`, 로컬 레지스트리

> 5번에서 만든 경량 카운터 이미지를 다른 사람도 받아 쓸 수 있게 레지스트리에 올립니다. Docker Hub 계정 없이, 로컬 레지스트리 컨테이너로 실습합니다.

## 상황

내 PC에만 있는 `counter:5.0` 이미지를 팀이 공유하려면 레지스트리에 올려야 합니다. 공식 저장소(`library/...`)에는 개인이 push할 수 없으니, **로컬 레지스트리**를 띄워 연습합니다.

> 주어진 코드는 없습니다. 앞 문제에서 만든 `counter:5.0` 이미지를 사용합니다. (없으면 5번 이미지를 먼저 빌드)

## 요구사항

1. 공식 레지스트리 이미지 `registry:2` 로 로컬 레지스트리를 **5001 포트**에 띄운다
2. `counter:5.0` 을 로컬 레지스트리용 이름 `localhost:5001/counter:5.0` 으로 **tag** 한다
3. 그 이름으로 **push** 한다
4. 로컬 이미지를 지운 뒤, 레지스트리에서 다시 **pull** 해서 동작을 확인한다

## 검증

```bash
# 레지스트리 기동
docker run -d -p 5001:5000 --name registry registry:2

# tag → push
docker tag counter:5.0 localhost:5001/counter:5.0
docker push localhost:5001/counter:5.0

# 로컬에서 지우고, 레지스트리에서 다시 받기
docker rmi counter:5.0 localhost:5001/counter:5.0
docker pull localhost:5001/counter:5.0
docker run --rm localhost:5001/counter:5.0 &  # 또는 -d -p 로 띄워 curl
```
