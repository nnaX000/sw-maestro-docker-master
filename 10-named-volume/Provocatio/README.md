# 문제 10. named volume 영속화 — 사라진 카운터 숫자 살리기

**내용 범위:** 컨테이너 휘발성, named volume, 데이터 영속성, 볼륨 마운트(`-v`)

> 1번에서 컨테이너를 지우면 방문 수가 0으로 초기화됐던 것 기억하시죠? 드디어 그 숫자를 살려둡니다. 이번 카운터는 숫자를 `/data/count` 파일에 저장합니다.

## 상황

이번 `main.go`는 방문 수를 메모리가 아니라 `/data/count` 파일에 기록하고, 시작할 때 그 파일에서 복원합니다. 하지만 그 파일이 **컨테이너 안(쓰기 레이어)** 에만 있으면 컨테이너를 지울 때 같이 사라집니다. **named volume**으로 `/data`를 컨테이너 바깥에 둬야 살아남습니다.

## 주어지는 것

```
Provocatio/app/   (main.go, go.mod)   ← /data/count 에 방문 수 저장
```

## 요구사항

1. `counter:10.0` 으로 빌드한다 (Dockerfile에서 `/data` 디렉터리를 미리 만들어 둘 것)
2. named volume `countdata` 를 만들어 `/data`에 마운트해 컨테이너를 띄운다 (`-v countdata:/data`)
3. 방문 수를 몇 번 올린 뒤, 컨테이너를 **삭제하고 재생성**해도 숫자가 유지되는지 확인한다
4. (대조) 볼륨 없이 띄우면 삭제 시 숫자가 0으로 돌아가는 것도 확인한다

## 검증

```bash
docker volume create countdata
docker run -d --name {{my-container}} -p 8080:5000 -v countdata:/data counter:10.0

curl http://localhost:8080   # visit #1
curl http://localhost:8080   # visit #2
curl http://localhost:8080   # visit #3

# 컨테이너를 통째로 지우고 다시 만들기 (같은 볼륨)
docker rm -f {{my-container}}
docker run -d --name {{my-container}} -p 8080:5000 -v countdata:/data counter:10.0

curl http://localhost:8080   # visit #4  ← 1~3이 살아남았다!
```
