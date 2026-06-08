# 해설 — 문제 12

## 1. 백업

```bash
docker run --rm \
  -v countdata:/data \
  -v "$(pwd)":/backup \
  alpine tar czf /backup/countdata.tar.gz -C /data .
```

읽어볼 포인트:
- `--rm` : 작업용 임시 컨테이너, 끝나면 자동 삭제.
- `-v countdata:/data` : 백업할 볼륨을 `/data`로 마운트.
- `-v "$(pwd)":/backup` : 결과 tar를 받을 호스트 현재 폴더를 `/backup`으로 bind mount.
- `tar czf ... -C /data .` : `/data` 안의 내용을 압축.

## 2. 복원

```bash
docker volume create countdata-restored
docker run --rm \
  -v countdata-restored:/data \
  -v "$(pwd)":/backup \
  alpine tar xzf /backup/countdata.tar.gz -C /data
```

## 3. 확인

```bash
docker run -d --name {{my-container}} -p 8080:5000 -v countdata-restored:/data counter:10.0
curl http://localhost:8080   # 백업 시점의 숫자 다음부터 이어짐
```

## 4. 패턴의 핵심

볼륨은 호스트 경로가 추상화돼 있어 직접 다루기 번거롭습니다. 그래서 **임시 컨테이너에 볼륨과 호스트 폴더를 함께 마운트하고 `tar`로 주고받는 것**이 표준 백업/이관 방법입니다.

> 주의: 실행 중인 **DB**의 데이터 볼륨을 그대로 tar하면 일관성이 깨질 수 있습니다. DB는 가능하면 엔진 덤프(`pg_dump`, `mysqldump`)를 병행하는 것이 안전합니다.
