package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func jobDir(id string) string {
	return filepath.Join(os.TempDir(), id)
}

func readJobResult(id string) (DetectResult, error) {
	data, err := os.ReadFile(filepath.Join(jobDir(id), "result.json"))
	if err != nil {
		return DetectResult{}, fmt.Errorf("job not found")
	}
	var result DetectResult
	if err := json.Unmarshal(data, &result); err != nil {
		return DetectResult{}, err
	}
	return result, nil
}

func readJobImage(id string) ([]byte, string, error) {
	matches, err := filepath.Glob(filepath.Join(jobDir(id), "annotated.*"))
	if err != nil || len(matches) == 0 {
		return nil, "", fmt.Errorf("annotated image not found")
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		return nil, "", err
	}
	contentType := http.DetectContentType(data)
	return data, contentType, nil
}
