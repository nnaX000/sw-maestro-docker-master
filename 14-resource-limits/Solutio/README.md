# 해설 — 문제 14

## 1. 실행과 관찰

```bash
docker build -t counter:14.0 ./app

# CPU 제한 — 0.5코어
docker run -d --name {{my-container}} -p 8080:5000 --cpus="0.5" counter:14.0
curl "http://localhost:8080/cpu?sec=15" &
docker stats --no-stream {{my-container}}   # CPU% ~50%에서 묶임

# 메모리 제한 — 64MB, 200MB 할당 시도 → OOM
docker rm -f {{my-container}}
docker run -d --name {{my-container}} -p 8080:5000 --memory="64m" counter:14.0
curl "http://localhost:8080/mem?mb=200"
docker inspect {{my-container}} --format '{{.State.OOMKilled}}'   # true
```

## 2. cgroups 란

**cgroups(Control Groups)** 는 프로세스 그룹이 쓸 수 있는 CPU·메모리·I/O 자원을 **제한·측정**하는 리눅스 커널 기능입니다. 컨테이너 격리의 한 축으로(다른 축은 namespace), Docker는 `--cpus`, `--memory` 같은 플래그를 이 cgroups 설정으로 변환합니다.

| 플래그 | 의미 |
|--------|------|
| `--cpus="0.5"` | CPU 0.5코어 분량으로 제한 (초과분은 throttle) |
| `--memory="64m"` | 메모리 상한 64MB (초과 시 OOM kill) |
| `--cpu-shares` | 상대적 CPU 우선순위 (경합 시) |
| `--memory-swap` | 메모리+스왑 합계 상한 |

## 3. 무슨 일이 일어났나

- **CPU**: `/cpu`는 모든 코어를 바쁘게 돌리려 하지만, cgroups가 0.5코어로 **throttle** 합니다. `docker stats`의 CPU%가 ~50%에서 더 못 올라갑니다. (죽지는 않고 느려짐)
- **메모리**: `/mem?mb=200`은 약 200MB를 점유하려 하지만 한도는 64MB입니다. 커널이 한도 초과를 감지하면 컨테이너 프로세스를 **OOM kill** 합니다. 종료 코드는 **137**(=128+SIGKILL 9), `inspect`의 `OOMKilled`는 `true`.

## 4. 실무 의미

- 한 컨테이너가 호스트 자원을 독식해 다른 서비스를 굶기는 일을 막으려면 **항상 리소스 한도**를 거는 것이 안전합니다.
- 메모리 한도는 앱의 실제 사용량보다 여유 있게 잡되, OOM이 잦다면 한도/앱 양쪽을 점검합니다.
- Kubernetes의 requests/limits도 결국 이 cgroups 위에서 동작합니다.
