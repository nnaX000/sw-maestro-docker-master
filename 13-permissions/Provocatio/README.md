# 문제 13. 권한(UID/GID) 트러블슈팅 — permission denied 잡기

**내용 범위:** 컨테이너 프로세스 UID, bind mount 권한 문제, `--user`, `chown`, named volume 회피

> 보안을 위해 컨테이너를 **비root 사용자**로 실행하면, bind mount한 호스트 디렉터리에 쓸 때 `permission denied`가 나는 일이 잦습니다. 재현하고 해결합니다.

## 상황

이번 이미지는 비root 사용자(uid 1000)로 실행됩니다. 카운터는 `/data/count`에 써야 하는데, 호스트의 root 소유 디렉터리를 bind mount하면 쓰기가 막힙니다.

## 주어지는 것

```
Provocatio/app/   (main.go, go.mod)   ← /data/count 에 쓰기 (실패 시 500 에러 반환)
```

Dockerfile은 비root 사용자로 실행하도록 작성해야 합니다(아래 요구사항).

## 요구사항

**(1) 문제 재현**
- Dockerfile에서 uid 1000 사용자를 만들고 `USER`로 지정해 `counter:13.0` 빌드
- root 소유의 호스트 디렉터리를 만들어 `/data`에 bind mount → 호출하면 쓰기 실패(`permission denied`)를 확인

**(2) 해결** — 아래 중 하나로 정상화
- 호스트 디렉터리를 `chown 1000:1000` 으로 소유자 변경, 또는
- `--user` 로 실행 UID를 호스트 디렉터리 소유자에 맞춤, 또는
- bind mount 대신 **named volume** 사용 (권한 자동 처리)

## 검증

```bash
# (1) 재현 — root 소유 호스트 디렉터리 bind mount
mkdir -p ~/hostdata   # 보통 root 또는 당신 uid 소유
docker run -d --name {{my-container}} -p 8080:5000 -v ~/hostdata:/data counter:13.0
curl http://localhost:8080
# → 쓰기 실패: ... permission denied

# (2-a) 해결: 호스트 디렉터리 소유자를 컨테이너 uid(1000)로
sudo chown -R 1000:1000 ~/hostdata
docker rm -f {{my-container}}
docker run -d --name {{my-container}} -p 8080:5000 -v ~/hostdata:/data counter:13.0
curl http://localhost:8080   # visit #1

# (2-b) 또는: named volume 으로 회피 (권한 자동 처리)
docker run -d --name {{my-container}} -p 8080:5000 -v countdata2:/data counter:13.0
```
