package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// // // // // // // // // //

func dirsWithInit(rootPath string) ([]string, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("rootPath помилка: %w", err)
	}
	if !info.IsDir() {
		return nil, errors.New("rootPath повинен бути директорією")
	}

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("не можу прочитати rootPath: %w", err)
	}

	var result []string

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		dirPath := filepath.Join(rootPath, e.Name())
		dirEntries, err := os.ReadDir(dirPath)
		if err != nil {
			continue
		}

		if len(dirEntries) == 0 {
			continue
		}

		hasInit := false
		for _, de := range dirEntries {
			if !de.IsDir() && de.Name() == "init.go" {
				hasInit = true
				break
			}
		}
		if !hasInit {
			continue
		}

		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			continue
		}
		result = append(result, absPath)
	}

	return result, nil
}
