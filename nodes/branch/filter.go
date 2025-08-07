package branch

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&FilterBranchNode{
		BaseNode: *node.NewBaseNode("teatime.branch.filter", node.NodeTypeBranch, "Filter", "데이터를 필터링하여 조건에 맞는 항목만 처리하는 필터 브랜치 노드입니다."),
	})
}

// 데이터를 필터링하여 분기 처리하는 필터 브랜치 노드
type FilterBranchNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *FilterBranchNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Filtered Items",
			Description: "필터링된 항목들",
			Key:         "filteredItems",
			Value:       "",
			Type:        node.JSONArray,
		},
		{
			Name:        "Rejected Items",
			Description: "필터링에서 제외된 항목들",
			Key:         "rejectedItems",
			Value:       "",
			Type:        node.JSONArray,
		},
		{
			Name:        "Filtered Count",
			Description: "필터링된 항목 수",
			Key:         "filteredCount",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Rejected Count",
			Description: "제외된 항목 수",
			Key:         "rejectedCount",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Total Count",
			Description: "전체 항목 수",
			Key:         "totalCount",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Filter Rate",
			Description: "필터링 비율 (%)",
			Key:         "filterRate",
			Value:       "",
			Type:        node.Float64,
		},
	}
}

func (r *FilterBranchNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Input Data",
			Description: "필터링할 입력 데이터 (배열)",
			Optional:    false,
			Key:         "inputData",
			Value:       "[]",
			Type:        node.JSON,
		},
		{
			Name:        "Filter Type",
			Description: "필터 타입",
			Optional:    false,
			Key:         "filterType",
			Value:       "simple",
			Type:        node.String,
			Options:     []string{"simple", "complex", "custom", "regex", "range", "unique", "duplicate"},
		},
		{
			Name:        "Filter Field",
			Description: "필터링할 필드명 (객체 배열의 경우)",
			Optional:    true,
			Key:         "filterField",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Filter Operator",
			Description: "필터 연산자",
			Optional:    true,
			Key:         "filterOperator",
			Value:       "==",
			Type:        node.String,
			Options:     []string{"==", "!=", ">", "<", ">=", "<=", "contains", "notContains", "startsWith", "endsWith", "in", "notIn", "between", "isNull", "isNotNull"},
		},
		{
			Name:        "Filter Value",
			Description: "필터 값",
			Optional:    true,
			Key:         "filterValue",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "Filter Expression",
			Description: "복잡한 필터 표현식 (custom 타입용)",
			Optional:    true,
			Key:         "filterExpression",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "Multiple Filters",
			Description: "다중 필터 조건 (JSON 배열)",
			Optional:    true,
			Key:         "multipleFilters",
			Value:       "[]",
			Type:        node.JSON,
		},
		{
			Name:        "Filter Logic",
			Description: "다중 필터 논리 연산",
			Optional:    true,
			Key:         "filterLogic",
			Value:       "AND",
			Type:        node.String,
			Options:     []string{"AND", "OR", "XOR", "CUSTOM"},
		},
		{
			Name:        "Case Sensitive",
			Description: "대소문자 구분",
			Optional:    true,
			Key:         "caseSensitive",
			Value:       "true",
			Type:        node.Bool,
		},
		{
			Name:        "Null Handling",
			Description: "NULL 값 처리 방식",
			Optional:    true,
			Key:         "nullHandling",
			Value:       "exclude",
			Type:        node.String,
			Options:     []string{"include", "exclude", "only"},
		},
		{
			Name:        "Sort After Filter",
			Description: "필터링 후 정렬 여부",
			Optional:    true,
			Key:         "sortAfterFilter",
			Value:       "false",
			Type:        node.Bool,
		},
		{
			Name:        "Sort Field",
			Description: "정렬 필드",
			Optional:    true,
			Key:         "sortField",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Sort Order",
			Description: "정렬 순서",
			Optional:    true,
			Key:         "sortOrder",
			Value:       "asc",
			Type:        node.String,
			Options:     []string{"asc", "desc"},
		},
		{
			Name:        "Limit",
			Description: "결과 제한 수",
			Optional:    true,
			Key:         "limit",
			Value:       "0",
			Type:        node.Int64,
		},
		{
			Name:        "Offset",
			Description: "시작 오프셋",
			Optional:    true,
			Key:         "offset",
			Value:       "0",
			Type:        node.Int64,
		},
	}
}

func (r *FilterBranchNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *FilterBranchNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
