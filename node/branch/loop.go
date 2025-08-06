package branch

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&LoopBranchNode{})
}

// 반복 처리를 위한 루프 브랜치 노드
type LoopBranchNode struct {
	customParams []types.NodeProperty
}

func (r *LoopBranchNode) Name() string {
	return "Loop"
}

func (r *LoopBranchNode) Description() string {
	return "조건에 따라 워크플로우를 반복 실행하는 루프 브랜치 노드입니다."
}

func (r *LoopBranchNode) Type() types.NodeType {
	return types.NodeTypeBranch
}

func (r *LoopBranchNode) ID() string {
	return "teatime.branch.loop"
}

func (r *LoopBranchNode) Output() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Current Index",
			Description: "현재 반복 인덱스",
			Key:         "currentIndex",
			Value:       "",
			Type:        types.Int64,
		},
		{
			Name:        "Current Item",
			Description: "현재 처리 중인 아이템",
			Key:         "currentItem",
			Value:       "",
			Type:        types.JSON,
			Optional:    true,
		},
		{
			Name:        "Loop Results",
			Description: "각 반복의 결과 배열",
			Key:         "loopResults",
			Value:       "",
			Type:        types.JSONArray,
		},
		{
			Name:        "Total Iterations",
			Description: "총 반복 횟수",
			Key:         "totalIterations",
			Value:       "",
			Type:        types.Int64,
		},
		{
			Name:        "Break Reason",
			Description: "루프 종료 이유",
			Key:         "breakReason",
			Value:       "",
			Type:        types.String,
			Optional:    true,
		},
	}
}

func (r *LoopBranchNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Loop Type",
			Description: "루프 타입",
			Optional:    false,
			Key:         "loopType",
			Value:       "for",
			Type:        types.String,
			Options:     []string{"for", "forEach", "while", "doWhile", "map", "filter", "reduce"},
		},
		{
			Name:        "Input Array",
			Description: "반복할 배열 (forEach, map, filter 타입용)",
			Optional:    true,
			Key:         "inputArray",
			Value:       "[]",
			Type:        types.JSON,
		},
		{
			Name:        "Start Index",
			Description: "시작 인덱스 (for 타입용)",
			Optional:    true,
			Key:         "startIndex",
			Value:       "0",
			Type:        types.Int64,
		},
		{
			Name:        "End Index",
			Description: "종료 인덱스 (for 타입용)",
			Optional:    true,
			Key:         "endIndex",
			Value:       "10",
			Type:        types.Int64,
		},
		{
			Name:        "Step",
			Description: "증가값 (for 타입용)",
			Optional:    true,
			Key:         "step",
			Value:       "1",
			Type:        types.Int64,
		},
		{
			Name:        "Condition",
			Description: "반복 조건 (while, doWhile 타입용)",
			Optional:    true,
			Key:         "condition",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Max Iterations",
			Description: "최대 반복 횟수 (무한 루프 방지)",
			Optional:    true,
			Key:         "maxIterations",
			Value:       "1000",
			Type:        types.Int64,
		},
		{
			Name:        "Break Condition",
			Description: "루프 중단 조건",
			Optional:    true,
			Key:         "breakCondition",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Continue Condition",
			Description: "현재 반복 건너뛰기 조건",
			Optional:    true,
			Key:         "continueCondition",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Parallel Execution",
			Description: "병렬 실행 여부",
			Optional:    true,
			Key:         "parallel",
			Value:       "false",
			Type:        types.Bool,
		},
		{
			Name:        "Batch Size",
			Description: "병렬 실행 시 배치 크기",
			Optional:    true,
			Key:         "batchSize",
			Value:       "10",
			Type:        types.Int64,
		},
		{
			Name:        "Accumulator",
			Description: "누적값 초기값 (reduce 타입용)",
			Optional:    true,
			Key:         "accumulator",
			Value:       "",
			Type:        types.JSON,
		},
	}
}

func (r *LoopBranchNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *LoopBranchNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
