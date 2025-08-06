package types

import (
	"context"
)

type NodeType string

const (
	NodeTypeTrigger NodeType = "trigger"
	NodeTypeBranch  NodeType = "branch"
	NodeTypeAction  NodeType = "action"
	NodeTypeUtil    NodeType = "util"
)

type BaseNode struct {
	ID          string
	Type        NodeType
	Description string
	Properties  []NodeProperty
}

// Node는 모든 노드가 구현해야 하는 기본 인터페이스
type Node interface {
	// 노드의 고유 식별자
	ID() string

	// 노드 이름
	Name() string

	// 노드 타입 반환
	Type() NodeType

	// 노드 설명
	Description() string

	// 노드의 입력 속성 정의
	Properties() []NodeProperty

	// 노드의 출력 속성 정의
	Output() []NodeProperty
}

// TriggerNode는 워크플로우를 시작하는 트리거 노드 인터페이스
type TriggerNode interface {
	Node

	// 트리거 활성화
	Activate(ctx context.Context) error

	// 트리거 비활성화
	Deactivate(ctx context.Context) error

	// 트리거 상태 확인
	IsActive() bool

	// 트리거 이벤트 수신 대기
	WaitForTrigger(ctx context.Context) (map[string]interface{}, error)
}

// ExecutableNode는 실행 가능한 노드를 위한 인터페이스
type ExecutableNode interface {
	Node

	// 노드 실행
	Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)

	// 노드 검증
	Validate(input map[string]interface{}) error

	// 실행 전 준비 작업
	Prepare(ctx context.Context) error

	// 실행 후 정리 작업
	Cleanup(ctx context.Context) error
}

// ActionNode는 작업을 수행하는 액션 노드 인터페이스
type ActionNode interface {
	ExecutableNode

	// 재시도 가능 여부
	IsRetryable() bool

	// 최대 재시도 횟수
	MaxRetries() int

	// 타임아웃 설정 (밀리초)
	Timeout() int64
}

// BranchNode는 분기 처리를 담당하는 브랜치 노드 인터페이스
type BranchNode interface {
	Node

	// 분기 조건 평가
	Evaluate(ctx context.Context, input map[string]interface{}) (string, error)

	// 가능한 분기 경로들
	GetBranches() []string

	// 기본 분기 경로
	DefaultBranch() string
}

// UtilNode는 유틸리티 기능을 제공하는 노드 인터페이스
type UtilNode interface {
	ExecutableNode

	// 동기/비동기 실행 여부
	IsAsync() bool

	// 캐시 가능 여부
	IsCacheable() bool

	// 캐시 키 생성
	CacheKey(input map[string]interface{}) string
}

// CustomizableNode는 커스텀 파라미터를 지원하는 노드 인터페이스
type CustomizableNode interface {
	Node

	// 커스텀 파라미터 반환
	CustomParams() []NodeProperty

	// 커스텀 파라미터 추가
	AddCustomParam(param NodeProperty)

	// 커스텀 파라미터 제거
	RemoveCustomParam(key string)

	// 커스텀 파라미터 초기화
	ClearCustomParams()
}

// NodeMetadata는 노드의 메타데이터를 정의하는 인터페이스
type NodeMetadata interface {
	// 노드 이름
	Name() string

	// 노드 버전
	Version() string

	// 노드 작성자
	Author() string

	// 노드 카테고리
	Category() string

	// 노드 태그들
	Tags() []string

	// 노드 아이콘
	Icon() string

	// 노드 색상
	Color() string

	// 노드 문서 URL
	DocumentationURL() string
}

// NodeStats는 노드 실행 통계를 위한 인터페이스
type NodeStats interface {
	// 실행 횟수
	ExecutionCount() int64

	// 성공 횟수
	SuccessCount() int64

	// 실패 횟수
	FailureCount() int64

	// 평균 실행 시간 (밀리초)
	AverageExecutionTime() int64

	// 마지막 실행 시간
	LastExecutionTime() int64

	// 통계 초기화
	ResetStats()
}

// NodeEventHandler는 노드 이벤트 처리를 위한 인터페이스
type NodeEventHandler interface {
	// 실행 시작 이벤트
	OnExecutionStart(ctx context.Context, input map[string]interface{})

	// 실행 완료 이벤트
	OnExecutionComplete(ctx context.Context, output map[string]interface{})

	// 실행 실패 이벤트
	OnExecutionError(ctx context.Context, err error)

	// 상태 변경 이벤트
	OnStateChange(oldState, newState string)
}

// NodeLifecycle는 노드 생명주기 관리를 위한 인터페이스
type NodeLifecycle interface {
	// 노드 초기화
	Initialize(config map[string]interface{}) error

	// 노드 시작
	Start(ctx context.Context) error

	// 노드 중지
	Stop(ctx context.Context) error

	// 노드 재시작
	Restart(ctx context.Context) error

	// 노드 상태 확인
	HealthCheck(ctx context.Context) error

	// 노드 파괴
	Destroy() error
}
