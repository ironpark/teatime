package branch

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&LoopBranchNode{
		BaseNode: node.NewBaseNode(
			"teatime.branch.loop",
			node.NodeTypeBranch,
			"Loop",
			"조건에 따라 워크플로우를 반복 실행하는 루프 브랜치 노드입니다.",
			"Repeat",
			[]node.NodeProperty{
				node.SelectProp("loopType", "Loop Type", []string{"for", "forEach", "while", "doWhile", "map", "filter", "reduce"},
					node.WithDescription("루프 타입"),
					node.RequiredWithDefault("for"),
				),
				node.JSONProp("inputArray", "Input Array",
					node.WithDescription("반복할 배열 (forEach, map, filter 타입용)"),
					node.OptionalWithDefault("[]"),
				),
				node.IntProp("startIndex", "Start Index",
					node.WithDescription("시작 인덱스 (for 타입용)"),
					node.OptionalWithDefault(int64(0)),
				),
				node.IntProp("endIndex", "End Index",
					node.WithDescription("종료 인덱스 (for 타입용)"),
					node.OptionalWithDefault(int64(10)),
				),
				node.IntProp("step", "Step",
					node.WithDescription("증가값 (for 타입용)"),
					node.OptionalWithDefault(int64(1)),
				),
				node.StringProp("condition", "Condition",
					node.WithDescription("반복 조건 (while, doWhile 타입용)"),
					node.Optional(),
				),
				node.IntProp("maxIterations", "Max Iterations",
					node.WithDescription("최대 반복 횟수 (무한 루프 방지)"),
					node.OptionalWithDefault(int64(1000)),
				),
				node.StringProp("breakCondition", "Break Condition",
					node.WithDescription("루프 중단 조건"),
					node.Optional(),
				),
				node.StringProp("continueCondition", "Continue Condition",
					node.WithDescription("현재 반복 건너뛰기 조건"),
					node.Optional(),
				),
				node.BoolProp("parallel", "Parallel Execution",
					node.WithDescription("병렬 실행 여부"),
					node.OptionalWithDefault(false),
				),
				node.IntProp("batchSize", "Batch Size",
					node.WithDescription("병렬 실행 시 배치 크기"),
					node.OptionalWithDefault(int64(10)),
				),
				node.JSONProp("accumulator", "Accumulator",
					node.WithDescription("누적값 초기값 (reduce 타입용)"),
					node.Optional(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.Int64, "currentIndex", "Current Index",
					node.WithDescription("현재 반복 인덱스"),
				),
				node.OutputProp(node.JSON, "currentItem", "Current Item",
					node.WithDescription("현재 처리 중인 아이템"),
				),
				node.OutputProp(node.JSON, "loopResults", "Loop Results",
					node.WithDescription("각 반복의 결과 배열"),
				),
				node.OutputProp(node.Int64, "totalIterations", "Total Iterations",
					node.WithDescription("총 반복 횟수"),
				),
				node.OutputProp(node.String, "breakReason", "Break Reason",
					node.WithDescription("루프 종료 이유"),
				),
			},
			nil, // Use default output handle
		),
	})
}

// 반복 처리를 위한 루프 브랜치 노드
type LoopBranchNode struct {
	node.BaseNode
}
