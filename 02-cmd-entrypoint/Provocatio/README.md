# 문제 02. 실행 인자 설계 — CMD vs ENTRYPOINT

**내용 범위:** `ENTRYPOINT`와 `CMD`의 차이와 조합, exec form, `docker run` 인자 덮어쓰기

> 1번 카운터에 인사 메시지 옵션(`-msg`)이 생겼습니다. 메시지를 바꿀 때마다 이미지를 새로 빌드하긴 번거롭습니다. 실행할 때 인자로 바꿀 수 있게 만들어 봅시다.

## 상황

운영팀이 같은 카운터 이미지를 환경마다 다른 인사말로 띄우고 싶어 합니다.

> "개발에선 'Hello (dev)', 운영에선 'Hello (prod)'로 띄우고 싶은데, 매번 다시 빌드하긴 싫어요. 실행할 때 메시지만 바꿔 끼울 수 없나요?"

`ENTRYPOINT`로 실행 파일은 고정하고, `CMD`로 기본 인자를 주되 `docker run`에서 덮어쓸 수 있게 설계하세요.

## 주어지는 것

```
Provocatio/
└── app/
    ├── main.go        ← 수정 금지 (-msg 플래그로 메시지를 받음)
    └── go.mod
```

`main.go`는 `-msg` 플래그로 인사 메시지를 받습니다. 안 주면 기본값 `Hello`를 씁니다.

## 요구사항

`app/`에 `Dockerfile`을 작성하세요.

1. 베이스는 `golang:1.22`, 작업 디렉터리는 `/app`
2. 빌드는 1번과 동일하게 (`go build -o server .`)
3. **`ENTRYPOINT`** 로 `./server`를 **항상 실행**되게 고정한다 (exec form 사용)
4. **`CMD`** 로 기본 인자 `-msg "Hello from Container!"` 를 제공한다 (exec form)
5. 포트 5000 명시, 이미지 이름 `counter:2.0`

## 검증

```bash
# 1) 기본 CMD 인자가 적용되는가
docker run -d --name {{my-container}} -p 8080:5000 counter:2.0
curl http://localhost:8080
# → Hello from Container! (visit #1)

# 2) run에 인자를 주면 CMD가 덮어써지는가 (ENTRYPOINT는 그대로)
docker rm -f {{my-container}}
docker run -d --name {{my-container}} -p 8080:5000 counter:2.0 -msg "Hello (prod)"
curl http://localhost:8080
# → Hello (prod) (visit #1)

# 3) ENTRYPOINT 자체를 임시로 바꿔보기 (디버깅용)
docker run --rm --entrypoint sh counter:2.0 -c "echo entrypoint overridden"
# → entrypoint overridden
```
