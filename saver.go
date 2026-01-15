package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	dirPerm  = 0755
	filePerm = 0644
)

// Saver сохраняет данные по URL в локальную файловую систему
type Saver struct {
	outputDir string
}

// NewSaver создаёт новый Saver с указанной выходной директорией
func NewSaver(outputDir string) (*Saver, error) {
	// Приведём путь к абсолютному и нормализуем
	absDir, err := filepath.Abs(outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", outputDir, err)
	}
	return &Saver{outputDir: absDir}, nil
}

// Save сохраняет данные по заданному URL
func (s *Saver) Save(u *url.URL, data []byte) (string, error) {
	localPath, err := s.urlToLocalPath(u)
	if err != nil {
		return "", fmt.Errorf("failed to convert URL to local path: %w", err)
	}

	fullPath := filepath.Join(s.outputDir, localPath)

	// Создаём все подкаталоги
	dir := filepath.Dir(fullPath)
	if err = os.MkdirAll(dir, dirPerm); err != nil {
		return "", fmt.Errorf("failed to create directories for %s: %w", fullPath, err)
	}

	// Записываем файл
	if err = os.WriteFile(fullPath, data, filePerm); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return fullPath, nil
}

// urlToLocalPath преобразует URL в относительный путь внутри зеркала
func (s *Saver) urlToLocalPath(u *url.URL) (string, error) {
	if u.Host == "" {
		return "", fmt.Errorf("URL has no host: %s", u.String())
	}

	// Начинаем с хоста
	p := u.Host

	// Добавляем путь
	cleanPath := strings.TrimPrefix(u.Path, "/")
	if cleanPath == "" {
		cleanPath = "index.html"
	} else {
		// Если путь заканчивается на '/', считаем, что это каталог → index.html
		if strings.HasSuffix(u.Path, "/") {
			cleanPath = path.Join(cleanPath, "index.html")
		}
	}

	p = path.Join(p, cleanPath)

	// Защита от path traversal (очень важно!)
	finalPath := filepath.FromSlash(p) // превращает / в \ на Windows
	if strings.Contains(finalPath, "..") {
		return "", fmt.Errorf("invalid path with '..': %s", finalPath)
	}

	return finalPath, nil
}
