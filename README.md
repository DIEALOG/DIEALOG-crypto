# DIEALOG-crypto

Go 기반 **비밀번호 및 기밀 정보 암호화 vault** 백엔드 애플리케이션입니다.

## 핵심 기능

### 1. 암호화된 자격 증명 저장소 (Encrypted Credential Vault)
- **기능**: 사용자별 계정 이름과 비밀번호를 암호화하여 저장 및 조회
- **암호화 방식**: 
  - 데이터 키(DEK)를 AWS KMS 방식 생성 및 관리
  - 마스터 키(KEK)로 DEK 암호화 후 Credential에 저장
  - 실제 비밀번호는 DEK로 AES-256-GCM 암호화
- **지원 연산**:
  - `StorePassword`: 신규 자격 증명 저장 (존재하면 업데이트)
  - `GetPassword`: 계정 이름 기반 복호화된 비밀번호 조회 (조회 즉시 메모리에서 파기)
  - `DeletePassword`: 계정 이름 기반 자격 증명 삭제

### 2. 내부 서비스 인증 (Internal Service Auth)
- **기능**: 모든 Credential API 라우트에 대해 `X-Internal-Secret-Key` 헤더 검증
- **보안**: 상수 시간 비교(Constant-Time comparison)로 타이밍 공격 방지
- **설정**: 환경 변수 `INTERNAL_SECRET_KEY`에 설정된 값과 헤더 값이 일치해야 접근 허용

### 3. KMS 키 관리 (KMS Integration)
- **기능**: 마스터 키를 통한 데이터 키 생성 및 암호화/복호화
- **구현**: 로컬 모의 KMS (AES-256-GCM 기반)
- **환경 설정**: `MOCK_KMS_MASTER_KEY` (32바이트 필수) 로컬 마스터 키 설정 가능

## 기술 스택

| 카테고리 | 기술 |
|----------|------|
| 언어 | Go 1.22+ |
| 웹 프레임워크 | Fiber v3 |
| ORM | Ent |
| 데이터베이스 | PostgreSQL |
| 암호화 | AES-256-GCM, UUID v4 |

## 아키텍처 흐름

1. **요청 진입**: Fiber 앱 → global middleware chain → `/api/credentials` 라우트
2. **인증 확인**: `X-Internal-Secret-Key` 헤더 검증 (미포함 / 불일치 시 401/403 반환)
3. **비즈니스 로직**: Use Case → Repository (Ent) → PostgreSQL
4. **암호화 연산**: KMS Service → DEK 생성/암호화/복호화 → AES-GCM 연산

## 디렉터리 구조

```
.
├── cmd/              # 메인 엔트리 포인트
├── api/              # 라우팅 및 핸들러
│   ├── handler/      # 요청 바인딩 및 응답 포맷팅
│   ├── middleware/   # 인증 및 공통 미들웨어
│   └── route/        # 라우트 등록 (Credential 라우트)
├── domain/           # 핵심 타입 및 인터페이스
├── infra/            # 구현체 (DB 레포, KMS)
│   ├── database/     # Ent 클라이언트 및 DSN
│   ├── kms/          # AWS KMS 연결부 (현재 Mock up 구현)
│   └── repository/   # Ent 기반 리포지토리
├── usecase/          # 비즈니스 로직 오케스트레이션
├── config/           # 환경 설정 (Viper)
└── ent/              # ORM 스키마 정의