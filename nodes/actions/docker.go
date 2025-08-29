package actions

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ironpark/teatime/internal/node"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/filters"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/client"
)

func init() {
	// Docker Container Run Node
	node.RegisterNode(&DockerRunNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.docker.run",
			node.NodeTypeAction,
			"Docker Run",
			"Docker 컨테이너를 실행합니다.",
			"Container",
			[]node.NodeProperty{
				node.StringProp("image", "Image",
					node.WithDescription("실행할 Docker 이미지 (예: nginx:latest)"),
					node.Required(),
				),
				node.StringProp("name", "Container Name",
					node.WithDescription("컨테이너 이름 (선택사항)"),
					node.Optional(),
				),
				node.StringProp("command", "Command",
					node.WithDescription("실행할 명령어 (선택사항)"),
					node.Optional(),
				),
				node.StringArrayProp("env", "Environment Variables",
					node.WithDescription("환경변수 목록 (KEY=VALUE 형식)"),
					node.Optional(),
				),
				node.StringArrayProp("ports", "Port Mappings",
					node.WithDescription("포트 매핑 목록 (8080:80 형식)"),
					node.Optional(),
				),
				node.BoolProp("detach", "Detach",
					node.WithDescription("백그라운드에서 실행"),
					node.OptionalWithDefault(true),
				),
				node.BoolProp("remove", "Auto Remove",
					node.WithDescription("컨테이너 종료 시 자동 삭제"),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "containerID", "Container ID",
					node.WithDescription("생성된 컨테이너 ID입니다."),
				),
				node.OutputProp(node.String, "logs", "Logs",
					node.WithDescription("컨테이너 실행 로그입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "컨테이너가 성공적으로 실행됨",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "컨테이너 실행 중 오류 발생",
				},
			},
		),
	})

	// Docker Image Pull Node
	node.RegisterNode(&DockerPullNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.docker.pull",
			node.NodeTypeAction,
			"Docker Pull",
			"Docker 이미지를 Pull합니다.",
			"Container",
			[]node.NodeProperty{
				node.StringProp("image", "Image",
					node.WithDescription("Pull할 Docker 이미지 (예: nginx:latest)"),
					node.Required(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "imageID", "Image ID",
					node.WithDescription("Pull된 이미지 ID입니다."),
				),
				node.OutputProp(node.String, "status", "Status",
					node.WithDescription("Pull 상태 메시지입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "이미지 Pull 성공",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "이미지 Pull 실패",
				},
			},
		),
	})

	// Docker Container Stop Node
	node.RegisterNode(&DockerStopNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.docker.stop",
			node.NodeTypeAction,
			"Docker Stop",
			"실행 중인 Docker 컨테이너를 중지합니다.",
			"Container",
			[]node.NodeProperty{
				node.StringProp("container", "Container",
					node.WithDescription("중지할 컨테이너 ID 또는 이름"),
					node.Required(),
				),
				node.IntProp("timeout", "Timeout",
					node.WithDescription("중지 대기 시간(초)"),
					node.OptionalWithDefault(int64(10)),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "status", "Status",
					node.WithDescription("컨테이너 상태입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "컨테이너 중지 성공",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "컨테이너 중지 실패",
				},
			},
		),
	})

	// Docker Container List Node
	node.RegisterNode(&DockerListNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.docker.list",
			node.NodeTypeAction,
			"Docker List",
			"Docker 컨테이너 목록을 조회합니다.",
			"Container",
			[]node.NodeProperty{
				node.BoolProp("all", "Show All",
					node.WithDescription("중지된 컨테이너도 포함"),
					node.OptionalWithDefault(false),
				),
				node.StringProp("status", "Status Filter",
					node.WithDescription("상태별 필터링 (created, restarting, running, removing, paused, exited, dead)"),
					node.Optional(),
				),
				node.StringProp("image", "Image Filter",
					node.WithDescription("이미지 이름으로 필터링"),
					node.Optional(),
				),
				node.StringProp("name", "Name Filter",
					node.WithDescription("컨테이너 이름으로 필터링"),
					node.Optional(),
				),
				node.StringArrayProp("labels", "Label Filters",
					node.WithDescription("레이블로 필터링 (key=value 형식)"),
					node.Optional(),
				),
				node.IntProp("limit", "Limit",
					node.WithDescription("반환할 컨테이너 개수 제한"),
					node.Optional(),
				),
				node.BoolProp("size", "Include Size",
					node.WithDescription("컨테이너 크기 정보 포함"),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.JSON, "containers", "Containers",
					node.WithDescription("컨테이너 목록입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "컨테이너 목록 조회 성공",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "컨테이너 목록 조회 실패",
				},
			},
		),
	})

	// Docker Container Remove Node
	node.RegisterNode(&DockerRemoveNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.docker.remove",
			node.NodeTypeAction,
			"Docker Remove",
			"Docker 컨테이너를 삭제합니다.",
			"Container",
			[]node.NodeProperty{
				node.StringProp("container", "Container",
					node.WithDescription("삭제할 컨테이너 ID 또는 이름"),
					node.Required(),
				),
				node.BoolProp("force", "Force Remove",
					node.WithDescription("강제 삭제 (실행 중인 컨테이너도 삭제)"),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "status", "Status",
					node.WithDescription("삭제 상태입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "컨테이너 삭제 성공",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "컨테이너 삭제 실패",
				},
			},
		),
	})
}

// DockerRunNode runs a Docker container
type DockerRunNode struct {
	node.BaseNode
}

func (d *DockerRunNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	imageName, ok := resolvedProps["image"].(string)
	if !ok || imageName == "" {
		return node.NodeResult{
			Error:         fmt.Errorf("image is required"),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"error": err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	defer cli.Close()

	// Container configuration
	config := &container.Config{
		Image: imageName,
	}

	// Set command if provided
	if command, ok := resolvedProps["command"].(string); ok && command != "" {
		config.Cmd = strings.Fields(command)
	}

	// Set environment variables
	if envVars, ok := resolvedProps["env"].([]string); ok {
		config.Env = envVars
	}

	// Host configuration
	hostConfig := &container.HostConfig{
		AutoRemove: false,
	}

	// Set auto remove
	if remove, ok := resolvedProps["remove"].(bool); ok {
		hostConfig.AutoRemove = remove
	}

	// Set port mappings
	if ports, ok := resolvedProps["ports"].([]string); ok && len(ports) > 0 {
		// Port mapping implementation would go here
		// This is simplified for brevity
	}

	// Set container name
	containerName := ""
	if name, ok := resolvedProps["name"].(string); ok && name != "" {
		containerName = name
	}

	// Create container
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"error": err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"containerID": resp.ID,
				"error":       err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Get logs if not detached
	logs := ""
	if detach, ok := resolvedProps["detach"].(bool); !ok || !detach {
		logReader, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err == nil {
			logBytes, _ := io.ReadAll(logReader)
			logs = string(logBytes)
			logReader.Close()
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"containerID": resp.ID,
			"logs":        logs,
			"error":       "",
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// DockerPullNode pulls a Docker image
type DockerPullNode struct {
	node.BaseNode
}

func (d *DockerPullNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	imageName, ok := resolvedProps["image"].(string)
	if !ok || imageName == "" {
		return node.NodeResult{
			Error:         fmt.Errorf("image is required"),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"error": err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	defer cli.Close()

	// Pull the image
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"error": err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	defer reader.Close()

	// Read pull status
	statusBytes, _ := io.ReadAll(reader)
	status := string(statusBytes)

	// Get image info
	imageInfo, err := cli.ImageInspect(ctx, imageName)
	imageID := ""
	if err == nil {
		imageID = imageInfo.ID
	}

	return node.NodeResult{
		Output: map[string]any{
			"imageID": imageID,
			"status":  status,
			"error":   "",
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// DockerStopNode stops a Docker container
type DockerStopNode struct {
	node.BaseNode
}

func (d *DockerStopNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	containerID, ok := resolvedProps["container"].(string)
	if !ok || containerID == "" {
		return node.NodeResult{
			Error:         fmt.Errorf("container is required"),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"error": err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	defer cli.Close()

	// Get timeout
	timeoutSecs := int64(10)
	if timeout, ok := resolvedProps["timeout"].(int64); ok && timeout > 0 {
		timeoutSecs = timeout
	}

	// Stop the container
	timeout := int(timeoutSecs)
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"status": "failed",
				"error":  err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"status": "stopped",
			"error":  "",
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// DockerListNode lists Docker containers
type DockerListNode struct {
	node.BaseNode
}

func (d *DockerListNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"error": err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	defer cli.Close()

	// Build list options from properties
	listOptions := container.ListOptions{}

	// Get all containers flag
	if showAll, ok := resolvedProps["all"].(bool); ok {
		listOptions.All = showAll
	}

	// Get size flag
	if includeSize, ok := resolvedProps["size"].(bool); ok {
		listOptions.Size = includeSize
	}

	// Get limit
	if limit, ok := resolvedProps["limit"].(int64); ok && limit > 0 {
		listOptions.Limit = int(limit)
	}

	// Build filters
	filterArgs := filters.NewArgs()

	// Status filter
	if status, ok := resolvedProps["status"].(string); ok && status != "" {
		filterArgs.Add("status", status)
	}

	// Image filter
	if imageFilter, ok := resolvedProps["image"].(string); ok && imageFilter != "" {
		filterArgs.Add("ancestor", imageFilter)
	}

	// Name filter
	if nameFilter, ok := resolvedProps["name"].(string); ok && nameFilter != "" {
		filterArgs.Add("name", nameFilter)
	}

	// Label filters
	if labels, ok := resolvedProps["labels"].([]string); ok && len(labels) > 0 {
		for _, label := range labels {
			filterArgs.Add("label", label)
		}
	}

	// Apply filters
	listOptions.Filters = filterArgs

	// List containers
	containers, err := cli.ContainerList(ctx, listOptions)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"error": err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Convert to output format
	containerList := make([]map[string]any, len(containers))
	for i, c := range containers {
		containerInfo := map[string]any{
			"id":      c.ID,
			"names":   c.Names,
			"image":   c.Image,
			"imageID": c.ImageID,
			"command": c.Command,
			"status":  c.Status,
			"state":   c.State,
			"ports":   c.Ports,
			"labels":  c.Labels,
			"created": c.Created,
		}

		// Include size information if requested
		if listOptions.Size {
			containerInfo["sizeRw"] = c.SizeRw
			containerInfo["sizeRootFs"] = c.SizeRootFs
		}

		containerList[i] = containerInfo
	}

	return node.NodeResult{
		Output: map[string]any{
			"containers": containerList,
			"error":      "",
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// DockerRemoveNode removes a Docker container
type DockerRemoveNode struct {
	node.BaseNode
}

func (d *DockerRemoveNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	containerID, ok := resolvedProps["container"].(string)
	if !ok || containerID == "" {
		return node.NodeResult{
			Error:         fmt.Errorf("container is required"),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"error": err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	defer cli.Close()

	// Get force flag
	force := false
	if forceRemove, ok := resolvedProps["force"].(bool); ok {
		force = forceRemove
	}

	// Remove the container
	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: force}); err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"status": "failed",
				"error":  err.Error(),
			},
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"status": "removed",
			"error":  "",
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}
