package actions

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&LsActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.ls",
			node.NodeTypeAction,
			"ls",
			"파일 및 하위 디렉토리 목록을 조회합니다.",
			"FolderOpen",
			[]node.NodeProperty{
				node.StringProp("path", "Directory Path",
					node.WithDescription("조회할 디렉토리 경로를 입력하세요"),
					node.Required(),
				),
				node.BoolProp("recursive", "Recursive",
					node.WithDescription("하위 디렉토리를 포함하여 파일 목록을 출력할지 여부를 선택하세요"),
					node.OptionalWithDefault(false),
				),
				node.SelectProp("filterType", "Filter Type", []string{"all", "files", "directories"},
					node.WithDescription("필터링할 타입을 선택하세요"),
					node.OptionalWithDefault("all"),
				),
				node.StringProp("pattern", "Pattern",
					node.WithDescription("파일 이름 패턴 (glob 패턴 지원)"),
					node.OptionalWithDefault("*"),
				),
				node.BoolProp("showHidden", "Show Hidden",
					node.WithDescription("숨김 파일 표시 여부"),
					node.OptionalWithDefault(false),
				),
				node.SelectProp("sortBy", "Sort By", []string{"name", "size", "modified", "created"},
					node.WithDescription("정렬 기준"),
					node.OptionalWithDefault("name"),
				),
				node.SelectProp("sortOrder", "Sort Order", []string{"asc", "desc"},
					node.WithDescription("정렬 순서"),
					node.OptionalWithDefault("asc"),
				),
				node.IntProp("maxDepth", "Max Depth",
					node.WithDescription("최대 탐색 깊이 (재귀 탐색 시)"),
					node.OptionalWithDefault(int64(1)),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.StringArray, "files", "Files",
					node.WithDescription("파일 목록입니다."),
				),
				node.OutputProp(node.StringArray, "directories", "Directories",
					node.WithDescription("디렉토리 목록입니다."),
				),
				node.OutputProp(node.Int64, "totalCount", "Total Count",
					node.WithDescription("전체 항목 수입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "Success output handle",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Error output handle",
				},
			},
		),
	})
}

// Run executes directory listing with filtering, sorting, and pagination options.
func (c *LsActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract parameters
	path, _ := resolvedProps["path"].(string)
	if path == "" {
		path = "."
	}

	recursive, _ := resolvedProps["recursive"].(bool)
	filterType, _ := resolvedProps["filterType"].(string)
	if filterType == "" {
		filterType = "all"
	}

	pattern, _ := resolvedProps["pattern"].(string)
	if pattern == "" {
		pattern = "*"
	}

	showHidden, _ := resolvedProps["showHidden"].(bool)
	sortBy, _ := resolvedProps["sortBy"].(string)
	if sortBy == "" {
		sortBy = "name"
	}

	sortOrder, _ := resolvedProps["sortOrder"].(string)
	if sortOrder == "" {
		sortOrder = "asc"
	}

	maxDepth, _ := resolvedProps["maxDepth"].(int64)
	if maxDepth <= 0 {
		maxDepth = 1
	}

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return node.NodeResult{
			Output: map[string]any{
				"files":       []string{},
				"directories": []string{},
				"totalCount":  int64(0),
				"error":       fmt.Sprintf("path does not exist: %s", path),
			},
			Error:         nil,
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	var files []string
	var directories []string
	var allEntries []fs.DirEntry

	// Walk directory
	err := filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip root directory itself
		if currentPath == path {
			return nil
		}

		// Calculate current depth
		relPath, _ := filepath.Rel(path, currentPath)
		depth := int64(strings.Count(relPath, string(filepath.Separator))) + 1

		// Respect maxDepth for recursive mode
		if recursive && depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip if not recursive and not immediate child
		if !recursive && depth > 1 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip hidden files if not requested
		if !showHidden && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Apply pattern filter
		if pattern != "*" {
			matched, _ := filepath.Match(pattern, d.Name())
			if !matched {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Collect entries
		allEntries = append(allEntries, d)

		// Categorize by type
		if d.IsDir() {
			if filterType == "all" || filterType == "directories" {
				directories = append(directories, currentPath)
			}
		} else {
			if filterType == "all" || filterType == "files" {
				files = append(files, currentPath)
			}
		}

		return nil
	})

	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"files":       []string{},
				"directories": []string{},
				"totalCount":  int64(0),
				"error":       err.Error(),
			},
			Error:         nil,
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Sort results
	sortEntries := func(entries []string) {
		sort.Slice(entries, func(i, j int) bool {
			switch sortBy {
			case "size":
				infoI, _ := os.Stat(entries[i])
				infoJ, _ := os.Stat(entries[j])
				if sortOrder == "desc" {
					return infoI.Size() > infoJ.Size()
				}
				return infoI.Size() < infoJ.Size()
			case "modified":
				infoI, _ := os.Stat(entries[i])
				infoJ, _ := os.Stat(entries[j])
				if sortOrder == "desc" {
					return infoI.ModTime().After(infoJ.ModTime())
				}
				return infoI.ModTime().Before(infoJ.ModTime())
			default: // name
				if sortOrder == "desc" {
					return entries[i] > entries[j]
				}
				return entries[i] < entries[j]
			}
		})
	}

	sortEntries(files)
	sortEntries(directories)

	totalCount := int64(len(files) + len(directories))

	return node.NodeResult{
		Output: map[string]any{
			"files":       files,
			"directories": directories,
			"totalCount":  totalCount,
			"error":       "",
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

type LsActionNode struct {
	node.BaseNode
}
