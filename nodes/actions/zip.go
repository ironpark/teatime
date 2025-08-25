package actions

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	// Zip Compress Node
	node.RegisterNode(&ZipCompressActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.zip.compress",
			node.NodeTypeAction,
			"Zip Compress",
			"파일 또는 디렉토리를 ZIP 파일로 압축합니다.",
			"Archive",
			[]node.NodeProperty{
				node.StringProp("sourcePath", "Source Path",
					node.WithDescription("압축할 파일 또는 디렉토리의 경로"),
					node.WithPlaceholder("/path/to/source"),
					node.Required(),
				),
				node.StringProp("targetPath", "Target ZIP Path",
					node.WithDescription("생성할 ZIP 파일의 경로"),
					node.WithPlaceholder("/path/to/archive.zip"),
					node.Required(),
				),
				node.BoolProp("overwrite", "Overwrite",
					node.WithDescription("기존 ZIP 파일을 덮어쓸지 여부"),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "zipPath", "ZIP Path",
					node.WithDescription("생성된 ZIP 파일의 경로입니다."),
				),
				node.OutputProp(node.Int64, "fileCount", "File Count",
					node.WithDescription("압축된 파일의 개수입니다."),
				),
				node.OutputProp(node.Int64, "compressedSize", "Compressed Size",
					node.WithDescription("압축된 파일의 크기(바이트)입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "ZIP compression completed successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to compress ZIP",
				},
			},
		),
	})

	// Zip Extract Node
	node.RegisterNode(&ZipExtractActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.zip.extract",
			node.NodeTypeAction,
			"Zip Extract",
			"ZIP 파일을 압축 해제합니다.",
			"Archive",
			[]node.NodeProperty{
				node.StringProp("zipPath", "ZIP Path",
					node.WithDescription("압축 해제할 ZIP 파일의 경로"),
					node.WithPlaceholder("/path/to/archive.zip"),
					node.Required(),
				),
				node.StringProp("targetDir", "Target Directory",
					node.WithDescription("압축 해제할 대상 디렉토리"),
					node.WithPlaceholder("/path/to/extract"),
					node.Required(),
				),
				node.BoolProp("overwrite", "Overwrite",
					node.WithDescription("기존 파일을 덮어쓸지 여부"),
					node.OptionalWithDefault(false),
				),
				node.BoolProp("createDir", "Create Directory",
					node.WithDescription("대상 디렉토리가 없으면 생성할지 여부"),
					node.OptionalWithDefault(true),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "extractDir", "Extract Directory",
					node.WithDescription("압축 해제된 디렉토리의 경로입니다."),
				),
				node.OutputProp(node.Int64, "extractedCount", "Extracted Count",
					node.WithDescription("압축 해제된 파일의 개수입니다."),
				),
				node.OutputProp(node.StringArray, "extractedFiles", "Extracted Files",
					node.WithDescription("압축 해제된 파일들의 경로 목록입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "ZIP extraction completed successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to extract ZIP",
				},
			},
		),
	})
}

// ZipCompressActionNode compresses files/directories into ZIP format
type ZipCompressActionNode struct {
	node.BaseNode
}

type zipCompressProps struct {
	SourcePath string `mapstructure:"sourcePath"`
	TargetPath string `mapstructure:"targetPath"`
	Overwrite  bool   `mapstructure:"overwrite"`
}

func (z *ZipCompressActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	var props zipCompressProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Check if source exists
	sourcePath := filepath.Clean(props.SourcePath)
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return node.NodeResult{
			Error:         fmt.Errorf("source path does not exist: %s", sourcePath),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Check if target already exists
	targetPath := filepath.Clean(props.TargetPath)
	if !props.Overwrite {
		if _, err := os.Stat(targetPath); err == nil {
			return node.NodeResult{
				Error:         fmt.Errorf("target ZIP file already exists: %s", targetPath),
				Continue:      false,
				OutputHandles: []string{"error"},
			}
		}
	}

	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to create target directory: %w", err),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Create ZIP file
	zipFile, err := os.Create(targetPath)
	if err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to create ZIP file: %w", err),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	var fileCount int64
	var compressedSize int64

	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path for ZIP entry
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		// Skip the source directory itself if it's a directory
		if relPath == "." {
			return nil
		}

		// Create ZIP entry
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Use forward slashes for ZIP paths
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}

			fileCount++
		}

		return nil
	})

	if err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to compress files: %w", err),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Get final file size
	if stat, err := zipFile.Stat(); err == nil {
		compressedSize = stat.Size()
	}

	output := map[string]any{
		"zipPath":        targetPath,
		"fileCount":      fileCount,
		"compressedSize": compressedSize,
	}

	return node.NodeResult{
		Output:        output,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// ZipExtractActionNode extracts ZIP files
type ZipExtractActionNode struct {
	node.BaseNode
}

type zipExtractProps struct {
	ZipPath   string `mapstructure:"zipPath"`
	TargetDir string `mapstructure:"targetDir"`
	Overwrite bool   `mapstructure:"overwrite"`
	CreateDir bool   `mapstructure:"createDir"`
}

func (z *ZipExtractActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	var props zipExtractProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Check if ZIP file exists
	zipPath := filepath.Clean(props.ZipPath)
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return node.NodeResult{
			Error:         fmt.Errorf("ZIP file does not exist: %s", zipPath),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Prepare target directory
	targetDir := filepath.Clean(props.TargetDir)
	if props.CreateDir {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return node.NodeResult{
				Error:         fmt.Errorf("failed to create target directory: %w", err),
				Continue:      false,
				OutputHandles: []string{"error"},
			}
		}
	} else {
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			return node.NodeResult{
				Error:         fmt.Errorf("target directory does not exist: %s", targetDir),
				Continue:      false,
				OutputHandles: []string{"error"},
			}
		}
	}

	// Open ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to open ZIP file: %w", err),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}
	defer reader.Close()

	var extractedCount int64
	extractedFiles := []string{}

	// Extract files
	for _, file := range reader.File {
		// Sanitize file path to prevent zip slip
		extractPath := filepath.Join(targetDir, file.Name)
		if !strings.HasPrefix(extractPath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return node.NodeResult{
				Error:         fmt.Errorf("invalid file path in ZIP: %s", file.Name),
				Continue:      false,
				OutputHandles: []string{"error"},
			}
		}

		// Check if file already exists
		if !props.Overwrite {
			if _, err := os.Stat(extractPath); err == nil {
				return node.NodeResult{
					Error:         fmt.Errorf("file already exists: %s", extractPath),
					Continue:      false,
					OutputHandles: []string{"error"},
				}
			}
		}

		if file.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(extractPath, file.FileInfo().Mode()); err != nil {
				return node.NodeResult{
					Error:         fmt.Errorf("failed to create directory: %w", err),
					Continue:      false,
					OutputHandles: []string{"error"},
				}
			}
		} else {
			// Create parent directory
			if err := os.MkdirAll(filepath.Dir(extractPath), 0755); err != nil {
				return node.NodeResult{
					Error:         fmt.Errorf("failed to create parent directory: %w", err),
					Continue:      false,
					OutputHandles: []string{"error"},
				}
			}

			// Extract file
			if err := extractFile(file, extractPath); err != nil {
				return node.NodeResult{
					Error:         fmt.Errorf("failed to extract file %s: %w", file.Name, err),
					Continue:      false,
					OutputHandles: []string{"error"},
				}
			}

			extractedCount++
			extractedFiles = append(extractedFiles, extractPath)
		}
	}

	output := map[string]any{
		"extractDir":     targetDir,
		"extractedCount": extractedCount,
		"extractedFiles": extractedFiles,
	}

	return node.NodeResult{
		Output:        output,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

func extractFile(file *zip.File, destPath string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}