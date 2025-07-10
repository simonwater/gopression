package fileutil

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PathBuilder 构建文件路径
func GetTestPath(directory, fileName string) string {
	outputDir, _ := GetProjectRoot()
	return filepath.Join(outputDir, ".temp", directory, fileName)
}

func GetProjectRoot() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// DeserializeObject 从文件反序列化对象 (使用 gob 编码)
func DeserializeObject[T any](filePath string) (T, error) {
	var obj T
	file, err := os.Open(filePath)
	if err != nil {
		return obj, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&obj); err != nil {
		return obj, fmt.Errorf("解码失败: %w", err)
	}
	return obj, nil
}

// DeserializeJSON 从文件反序列化 JSON 对象
func DeserializeJSON[T any](filePath string) (T, error) {
	var obj T
	file, err := os.Open(filePath)
	if err != nil {
		return obj, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return obj, fmt.Errorf("读取文件失败: %w", err)
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return obj, fmt.Errorf("JSON 解析失败: %w", err)
	}
	return obj, nil
}

// SerializeObject 序列化对象到文件 (使用 gob 编码)
func SerializeObject(obj interface{}, filePath string) error {
	if err := CreateParentIfNotExist(filePath); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(obj); err != nil {
		return fmt.Errorf("编码失败: %w", err)
	}
	return nil
}

// SerializeJSON 序列化对象为 JSON 到文件
func SerializeJSON(obj interface{}, filePath string, indent bool) error {
	if err := CreateParentIfNotExist(filePath); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if indent {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(obj); err != nil {
		return fmt.Errorf("JSON 编码失败: %w", err)
	}
	return nil
}

// WriteString 写入字符串到文件
func WriteString(content, filePath string) error {
	if err := CreateParentIfNotExist(filePath); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	return nil
}

// WriteAllLines 写入多行文本到文件
func WriteAllLines(lines []string, filePath string) error {
	if err := CreateParentIfNotExist(filePath); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return fmt.Errorf("写入行失败: %w", err)
		}
	}
	return nil
}

// CreateParentIfNotExist 创建父目录（如果不存在）
func CreateParentIfNotExist(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}
	return nil
}
