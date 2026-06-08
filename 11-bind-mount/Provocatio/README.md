# 문제 11. bind mount — 설정을 호스트에서 갈아끼우기

**내용 범위:** bind mount, 호스트 경로 직접 마운트, Volume과의 용도 차이

> 이번 카운터는 인사말을 `/config/message.txt` 파일에서 읽습니다(매 요청마다). 이미지를 다시 빌드하지 않고, **호스트의 파일을 컨테이너에 직접 연결(bind mount)** 해서 인사말을 바꿔봅니다.

## 상황

- 데이터(`/data/count`)는 10번처럼 **named volume**에 둡니다.
- 설정(`/config/message.txt`)은 **bind mount**로 호스트의 파일을 직접 연결합니다.
- 호스트에서 그 파일을 수정하면, 컨테이너 안에서도 즉시 바뀐 내용이 보여야 합니다.

## 주어지는 것

```
Provocatio/app/   (main.go, go.mod)
  - /data/count        에 방문 수 저장 (volume)
  - /config/message.txt 에서 인사말을 읽음 (bind mount 대상)
```

## 요구사항

1. `counter:11.0` 으로 빌드 (`/data`, `/config` 디렉터리 미리 생성)
2. 호스트에 인사말 파일을 하나 만든다 (예: `~/myconfig/message.txt` 에 `Hello (bind!)`)
3. named volume(`countdata`)은 `/data`에, **호스트 파일을 `/config/message.txt`에 bind mount** 하여 띄운다
4. 호스트의 `message.txt`를 고친 뒤 다시 호출하면 인사말이 바뀌는지 확인

## 검증

```bash
mkdir -p ~/myconfig && echo "Hello (bind!)" > ~/myconfig/message.txt

docker run -d --name {{my-container}} -p 8080:5000 \
  -v countdata:/data \
  -v ~/myconfig/message.txt:/config/message.txt \
  counter:11.0

curl http://localhost:8080
# → Hello (bind!) (visit #N)

# 호스트에서 설정만 바꿔치기 (재빌드 X)
echo "안녕하세요 (수정됨)" > ~/myconfig/message.txt
curl http://localhost:8080
# → 안녕하세요 (수정됨) (visit #N+1)
```
