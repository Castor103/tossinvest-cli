# Toss Securities Open API 마이그레이션 계획

토스증권이 공식 Open API 사전 신청을 시작했습니다 (광고 심의 시작일 2026-05-14).
이 문서는 `tossctl` 이 공식 API 출시 흐름에 맞춰 어떻게 진화할지 정리한 living
document 입니다. issue [#31](https://github.com/JungHoonGhae/tossinvest-cli/issues/31)
이 트래킹 anchor.

> **status:** Phase 0 (사전 신청 진행 중) · 마지막 업데이트 2026-05-19
>
> **⚠ 본 문서의 계획·timeline·phase 정의·우선순위는 사전 공지 없이 바뀔 수 있습니다.**
> 토스 공식 표면이 실제로 드러나는 시점, 자원/시간 사정, 새로운 정보에 따라 유연하게
> 재조정합니다. 여기 적힌 내용은 commitment 가 아니라 "현 시점 stance" 입니다.

## 토스 Open API 의 윤곽 (corp.tossinvest.com/ko/open-api 기준)

아래는 마케팅 페이지에 노출된 example 코드와 본문에서 읽은 내용입니다. 실제 endpoint
스펙·헤더·필드는 토큰 발급 후 직접 확인하기 전까지 *추정* 입니다.

- **Base URL:** `https://openapi.tossinvest.com/v1` *(페이지 예시 코드 기준, 미검증)*
- **인증:** `Authorization: Bearer <token>` + `X-Tossinvest-Account: <accountSeq>` 헤더 *(페이지 예시 기준, 미검증)*
- **프로토콜:** REST + WebSocket *(페이지 본문)*
- **공개 표면:** 시세 (실시간 호가/체결/캔들), 주문 (국내+해외 통합 마케팅), 계좌 조회, 종목/시장 정보. 각 표면의 endpoint 와 거래 권한 모델은 미공개
- **자격:** 토스증권 계좌 보유자만 사전 신청 가능 *(페이지 본문)*
- **출시 일자:** 명시되지 않음. 사전 신청 후 순차 롤링은 일반 패턴 추정일 뿐, 토스가 어떤 순서/속도로 푸는지는 미공개

## 우리 포지셔닝 변화

현재: "토스증권 웹 세션을 reverse-engineer 한 비공식 CLI"

공식 API 출시 후: **"한국 증권사를 AI 에이전트에 통일된 인터페이스로 연결하는 CLI"** —
백엔드 plugin 으로 추상화해서 official Toss / 비공식 Toss / (장기적으로) KIS · 키움 등을
같은 명령어로 다룸. 사용자가 토큰을 받았든 안 받았든 `tossctl portfolio positions` 의
표면은 동일.

## Phase 별 계획 (잠정)

표 안의 phase 정의·동작·작업 모두 잠정. 공식 표면이 드러나면 항목이 합쳐지거나
분리되거나 순서가 바뀔 수 있습니다.

| Phase | 트리거 | tossctl 동작 | 작업 |
|---|---|---|---|
| **0** *(지금)* | 사전 신청만 가능, 토큰 발급 0 | 현행 — session-based 만 | issue #31 트래킹, 사전 신청, 본 문서 유지 |
| **1** | 일부 사용자 토큰 발급 시작 | `tossctl auth login --official-token <token>` 추가. config 에 토큰 있으면 official, 없으면 session. **명령어 표면 동일** | `Broker` interface 추상화 (`TossSessionBroker` / `TossOfficialBroker`), `OAuthBearer` 인증, doctor 안내 |
| **2** | 대부분 토큰 발급, official 안정 | default 가 official, session 은 fallback. doctor 가 자동 전환 권장 | 거래 권한 모델 정리 (official 의 trading scope 가 별도 신청이라면 분기) |
| **3** | 정착 | session-based deprecation. KIS/키움 broker plugin 검토 | `tossctl --broker toss|kis|kiwoom` 가능성 평가 |

## Phase 1 의 UX 원칙

- 토큰 받은 사용자: `tossctl auth login --official-token ...` 한 번. 이후 끝
- 토큰 못 받은 사용자: 기존 흐름 그대로. 아무것도 안 변함
- **두 그룹이 같은 README, 같은 명령어, 같은 AGENTS.md** 를 봄. tossctl 이 매개

doctor 출력 예시:
```
Backend: toss-session (active)
Official API: not yet (waitlist) — apply at https://corp.tossinvest.com/ko/open-api
```

토큰 발급 후 (정확한 필드명/만료 정책은 실제 토큰 받은 후 확정):
```
Backend: toss-official (token expires ...)
Session fallback: configured
```

## 위험 요소

1. **비공식 endpoint 차단 가능성** — 공식 출시 후 토스가 reverse-engineered 접근을
   정책/기술적으로 막을 수 있음. 그 시점에 session 백엔드는 빠르게 deprecate 가 강제될
   수 있음
2. **공식 API 의 거래 권한** — official 이 거래 권한을 별도 신청해야 할 가능성 (대부분
   증권사가 그럼). 이 경우 tossctl 의 거래 기능은 official 백엔드에서 분기 처리 필요
3. **추상화 over-engineering** — 공식 표면을 직접 보기 전에 interface 를 짜면 잘못
   잡을 확률 ↑. (관련 결정은 아래 *결정 log* 참조)

## 결정 log

각 항목은 *그 시점의* stance. 이후 정보로 뒤집힐 수 있고, 뒤집힐 때 사전 공지하지
않습니다 — 새 항목을 추가할 뿐입니다.

- **2026-05-19** — issue #31 등록 (제보: @DaeHyeoNi). 사전 신청 페이지 확인. 사전 신청 진행. tossctl 의 일반화 (multi-broker) 방향을 *현재로서는* 선호. Phase 1 진입 전까지 코드 추상화는 보류 — 공식 표면을 보기 전 추상화는 잘못 잡을 확률이 크다는 판단

## 외부 contributor / 사용자에게 부탁

1. 토스 Open API 토큰을 받으셨다면 issue #31 에 댓글로 알려주세요 — phase 1 진입
   판단의 가장 빠른 신호. 토큰/계정번호 같은 민감 정보는 절대 공개 댓글에 붙이지 마세요
2. 공식 API 의 endpoint/스펙 문서를 발견하시면 issue #31 에 링크 공유 환영
3. 이 문서의 timeline/계획에 의견 있으면 issue #31 댓글 또는 별도 PR 환영
