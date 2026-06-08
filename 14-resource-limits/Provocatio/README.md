# 문제 14. 리소스 제한 (cgroups) — 한도를 걸고 관찰하기

**내용 범위:** cgroups 기반 리소스 제한(`--memory`, `--cpus`), OOM kill, CPU throttling, `docker stats`

> 컨테이너 격리의 한 축인 cgroups를 직접 체감합니다. 카운터에 부하 엔드포인트(`/cpu`, `/mem`)가 추가됐습니다. 한도를 걸고 무슨 일이 생기는지 봅니다.

## 상황

`main.go`에는 부하용 엔드포인트가 있습니다.
- `GET /cpu?sec=N` : N초 동안 모든 코어를 바쁘게 돌림
- `GET /mem?mb=N` : 약 N MB를 실제로 할당해 붙잡음

이 컨테이너에 CPU·메모리 한도를 걸고 동작을 관찰하세요.

## 주어지는 것

```
Provocatio/app/   (main.go, go.mod)   ← /cpu, /mem 부하 엔드포인트
```

## 요구사항

`counter:14.0` 빌드 후:

1. **CPU 제한** — `--cpus="0.5"` 로 띄우고 `/cpu`를 호출, `docker stats`에서 CPU%가 약 50%로 제한(throttle)되는지 관찰
2. **메모리 제한** — `--memory="64m"` 으로 띄우고 `/mem?mb=200`을 호출 → 한도를 넘으면 컨테이너가 **OOM kill** 되는지 확인 (`docker inspect`의 OOMKilled)

## 검증

```bash
docker build -t counter:14.0 ./app

# (1) CPU 제한
docker run -d --name {{my-container}} -p 8080:5000 --cpus="0.5" counter:14.0
curl "http://localhost:8080/cpu?sec=15"
docker stats --no-stream {{my-container}}     # CPU % 가 ~50% 근처에서 묶임

# (2) 메모리 제한 → OOM
docker rm -f {{my-container}}
docker run -d --name {{my-container}} -p 8080:5000 --memory="64m" counter:14.0
curl "http://localhost:8080/mem?mb=200"        # 64m 한도 초과 시도
docker inspect {{my-container}} --format '{{.State.OOMKilled}}'
# → true  (커널이 메모리 초과로 죽임)
docker ps -a    # STATUS: Exited (137)  ← 137 = OOM kill 신호(128+9)
```
