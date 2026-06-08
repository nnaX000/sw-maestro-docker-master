# 문제 12. 볼륨 백업과 복원 — 데이터를 안전하게 옮기기

**내용 범위:** named volume 백업/복원, 임시 `--rm` 컨테이너 + `tar`, 데이터 이관

> 10번에서 만든 `countdata` 볼륨의 데이터를 파일로 백업하고, 새 볼륨에 복원해 봅니다. 별도 코드는 없습니다.

## 상황

볼륨은 컨테이너 바깥에 있지만, 그 자체를 직접 `cp` 하긴 번거롭습니다. **임시 컨테이너에 볼륨을 마운트해 `tar`로 묶어내는** 표준 패턴을 익힙니다.

> 준비물: 10번/11번에서 데이터가 쌓인 `countdata` 볼륨 (없으면 10번을 먼저 실행)

## 요구사항

1. `countdata` 볼륨을 `tar.gz` 로 **백업** (현재 디렉터리에 `countdata.tar.gz` 생성)
2. 새 볼륨 `countdata-restored` 를 만들고, 백업본을 **복원**
3. 복원한 볼륨으로 카운터를 띄워 숫자가 그대로인지 확인

## 검증

```bash
# 백업 — alpine 임시 컨테이너에 볼륨과 현재 폴더를 함께 마운트
docker run --rm -v countdata:/data -v "$(pwd)":/backup alpine \
  tar czf /backup/countdata.tar.gz -C /data .

ls -lh countdata.tar.gz   # 백업 파일 생성 확인

# 복원
docker volume create countdata-restored
docker run --rm -v countdata-restored:/data -v "$(pwd)":/backup alpine \
  tar xzf /backup/countdata.tar.gz -C /data

# 복원 볼륨으로 기동 — 숫자 유지 확인
docker run -d --name {{my-container}} -p 8080:5000 -v countdata-restored:/data counter:10.0
curl http://localhost:8080
