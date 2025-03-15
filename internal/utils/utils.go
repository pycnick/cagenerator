package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

// ReadFile читает файл и возвращает его содержимое
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}
	return data, nil
}

// GetModuleName читает имя модуля из go.mod файла
func GetModuleName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("opening go.mod file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(line[7:]), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("reading go.mod: %w", err)
	}

	return "", fmt.Errorf("module declaration not found in go.mod")
}

// FormatProject запускает go mod tidy и go fmt для проекта
func FormatProject() error {
	if err := runCommand("go", "get", "go.uber.org/mock/mockgen"); err != nil {
		return fmt.Errorf("running go mod tidy: %w", err)
	}

	if err := runCommand("go", "generate", "./..."); err != nil {
		return fmt.Errorf("running go generate: %w", err)
	}

	if err := runCommand("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("running go mod tidy: %w", err)
	}

	if err := runCommand("go", "fmt", "./..."); err != nil {
		return fmt.Errorf("running go fmt: %w", err)
	}

	return nil
}

// runCommand выполняет команду и возвращает ошибку
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command %s failed: %s: %w", name, string(output), err)
	}
	return nil
}

// CamelToSnake конвертирует строку из camelCase в snake_case
func CamelToSnake(s string) string {
	var result []rune
	for i, c := range s {
		if i > 0 && unicode.IsUpper(c) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(c))
	}
	return string(result)
}

// EnsureDirectory создает директорию, если она не существует
func EnsureDirectory(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory %s: %w", path, err)
	}
	return nil
}

// FileExists проверяет существование файла
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Lower конвертирует строку в нижний регистр
func Lower(s string) string {
	return strings.ToLower(s)
}
