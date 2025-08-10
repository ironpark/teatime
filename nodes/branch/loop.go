package branch

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&LoopBranchNode{
		BaseNode: *node.NewBaseNode("teatime.branch.loop", node.NodeTypeBranch, "Loop", "조건에 따라 워크플로우를 반복 실행하는 루프 브랜치 노드입니다.", "Repeat"),
	})
}

// 반복 처리를 위한 루프 브랜치 노드
type LoopBranchNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *LoopBranchNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Current Index",
			Description: "현재 반복 인덱스",
			Key:         "currentIndex",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Current Item",
			Description: "현재 처리 중인 아이템",
			Key:         "currentItem",
			Value:       "",
			Type:        node.JSON,
			Optional:    true,
		},
		{
			Name:        "Loop Results",
			Description: "각 반복의 결과 배열",
			Key:         "loopResults",
			Value:       "",
			Type:        node.JSONArray,
		},
		{
			Name:        "Total Iterations",
			Description: "총 반복 횟수",
			Key:         "totalIterations",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Break Reason",
			Description: "루프 종료 이유",
			Key:         "breakReason",
			Value:       "",
			Type:        node.String,
			Optional:    true,
		},
	}
}

func (r *LoopBranchNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Loop Type",
			Description: "루프 타입",
			Optional:    false,
			Key:         "loopType",
			Value:       "for",
			Type:        node.String,
			Options:     []string{"for", "forEach", "while", "doWhile", "map", "filter", "reduce"},
		},
		{
			Name:        "Input Array",
			Description: "반복할 배열 (forEach, map, filter 타입용)",
			Optional:    true,
			Key:         "inputArray",
			Value:       "[]",
			Type:        node.JSON,
		},
		{
			Name:        "Start Index",
			Description: "시작 인덱스 (for 타입용)",
			Optional:    true,
			Key:         "startIndex",
			Value:       "0",
			Type:        node.Int64,
		},
		{
			Name:        "End Index",
			Description: "종료 인덱스 (for 타입용)",
			Optional:    true,
			Key:         "endIndex",
			Value:       "10",
			Type:        node.Int64,
		},
		{
			Name:        "Step",
			Description: "증가값 (for 타입용)",
			Optional:    true,
			Key:         "step",
			Value:       "1",
			Type:        node.Int64,
		},
		{
			Name:        "Condition",
			Description: "반복 조건 (while, doWhile 타입용)",
			Optional:    true,
			Key:         "condition",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Max Iterations",
			Description: "최대 반복 횟수 (무한 루프 방지)",
			Optional:    true,
			Key:         "maxIterations",
			Value:       "1000",
			Type:        node.Int64,
		},
		{
			Name:        "Break Condition",
			Description: "루프 중단 조건",
			Optional:    true,
			Key:         "breakCondition",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Continue Condition",
			Description: "현재 반복 건너뛰기 조건",
			Optional:    true,
			Key:         "continueCondition",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Parallel Execution",
			Description: "병렬 실행 여부",
			Optional:    true,
			Key:         "parallel",
			Value:       "false",
			Type:        node.Bool,
		},
		{
			Name:        "Batch Size",
			Description: "병렬 실행 시 배치 크기",
			Optional:    true,
			Key:         "batchSize",
			Value:       "10",
			Type:        node.Int64,
		},
		{
			Name:        "Accumulator",
			Description: "누적값 초기값 (reduce 타입용)",
			Optional:    true,
			Key:         "accumulator",
			Value:       "",
			Type:        node.JSON,
		},
	}
}

func (r *LoopBranchNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *LoopBranchNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
