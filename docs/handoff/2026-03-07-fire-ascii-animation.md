# Fire ASCII Animation

## 브랜치

- 현재: `doyoon/donut-animation`

## 완료된 작업

- `internal/fire/fire.go` — Doom fire 알고리즘 기반 ASCII 불꽃 애니메이션 구현
- `cmd/fire.go` — `fire` 서브커맨드 진입점
- `main.go` — `fire` 서브커맨드 분기 추가
- 8개 화염 zone이 가로 화면을 빈틈 없이 채우도록 구현
- 가우시안 가중 보간으로 zone 간 높이를 부드럽게 연결 (경계 끊김 해결)
- 애니메이션 속도 50ms → 100ms로 조정

## 핵심 기술 결정

| 결정 | 이유 |
|------|------|
| 단일 통합 버퍼 | 독립 버퍼 방식은 화염 경계에서 세로 직선 끊김 발생. 통합 버퍼에서 열 전파가 경계를 자연스럽게 넘어감 |
| 가우시안 보간 높이 | `exp(-dist²/(segW²×0.8))`로 각 x좌표의 화염 높이를 부드럽게 결정. 8개 zone 높이가 매끄러운 곡선으로 연결 |
| 70자 ASCII 팔레트 | 문자 획 밀도 기반 그레이스케일. heatLevel(0~19)을 팔레트 인덱스로 매핑 |
| 사인파 바람 | 각 zone마다 다른 phase의 sin(frame×0.05+phase)×0.8로 비동기적 일렁임 |

## 현재 상태

- 빌드 성공, `./.tui fire`로 실행 가능
- 기본 동작 검증 완료 (firetest 스크립트로 TTY 없이 출력 확인)

## TODO

- [ ] 화염 상단 윤곽이 더 자연스럽게 보이도록 cutoff 페이드 파라미터 튜닝
- [ ] 실제 터미널에서 시각적 완성도 최종 확인 (사용자 피드백 반영 중이었음)
- [ ] donut 애니메이션과 fire 애니메이션 간 전환 또는 메뉴 시스템 고려
- [ ] 커밋 (현재 미커밋 상태 — cmd/fire.go, cmd/donut.go, internal/fire/, internal/donut/, main.go)

## 미커밋 파일

```
 M main.go
?? cmd/donut.go
?? cmd/fire.go
?? internal/donut/
?? internal/fire/
?? pinterest-reference-2.png
?? pinterest-reference-3.png
?? pinterest-reference.png
?? tui
```

## 주의사항

- `tui` 바이너리가 untracked 상태로 남아있음 (.gitignore 확인 필요)
- pinterest-reference 이미지들은 디자인 참조용, 커밋 포함 여부 사용자 판단 필요
- `.playwright-mcp/` 디렉토리도 untracked (gitignore 대상)

## 참고 파일 구조

```
cmd/
  donut.go          — donut 서브커맨드
  fire.go           — fire 서브커맨드
internal/
  donut/donut.go    — donut 애니메이션 (기존, 변경 없음)
  fire/fire.go      — fire 애니메이션 (이번 세션에서 작성)
main.go             — 서브커맨드 라우팅
```
