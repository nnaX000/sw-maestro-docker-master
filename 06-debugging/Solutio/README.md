# 해설 — 문제 06

## 1. 진단 과정

```bash
# (1) 상태 확인 — Up 이 아니라 Exited 다
docker ps -a
# STATUS: Exited (1) 3 seconds ago

# (2) 로그에서 원인 — 앱이 남긴 fatal 메시지
docker logs {{my-container}}
# FATAL: GREETING 환경변수가 비어 있습니다. -e GREETING=... 로 전달하세요

# (3) 종료 코드 확인
docker inspect {{my-container}} --format '{{.State.ExitCode}}'
# 1   ← 비정상 종료
```

원인: 앱은 시작할 때 `GREETING` 환경변수를 요구하는데, 그게 비어 있어 `log.Fatal`로 즉시 종료된 것입니다. `docker run`은 "컨테이너 생성"에 성공했으니 에러가 안 났지만, 안의 프로세스(PID 1)가 죽으면서 컨테이너도 함께 종료됐습니다.

## 2. 해결 — 환경변수 주입

```bash
docker rm -f {{my-container}}
docker run -d --name {{my-container}} -p 8080:5000 \
  -e GREETING="Hello from Container!" counter:6.0

curl http://localhost:8080
# → Hello from Container! (visit #1)
```

## 3. 진단 명령 정리

| 명령 | 용도 |
|------|------|
| `docker ps -a` | 죽은 컨테이너까지 상태 확인 (`-a` 없으면 Up만 보임) |
| `docker logs <c>` | PID 1의 STDOUT/STDERR — 죽은 이유의 1순위 단서 |
| `docker inspect <c>` | ExitCode, 재시작 횟수, 마운트, 네트워크 등 메타데이터 |
| `docker exec -it <c> sh` | **떠 있는** 컨테이너 내부 진입 (죽은 컨테이너엔 불가) |
| `docker cp <c>:/path ./` | 컨테이너 ↔ 호스트 파일 복사 (정지 상태에서도 가능) |

참고로 cp 명령어가 가능한 이유는 container 의 UnionFS 덕분이란거! 
모든 Container 기술은 결국 Host의 File System에 촘촘히 lowerdir, upperdir, mergedir 로 표현되는거 기억나시죠~?

## 4. 교훈

"`docker run`이 성공 = 서비스 정상"이 아닙니다. run은 컨테이너를 만들 뿐이고, 안의 프로세스가 살아 있어야 서비스가 됩니다. 컨테이너가 안 보이면 **`docker ps -a` → `docker logs`** 가 진단의 출발점입니다.
