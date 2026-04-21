package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

const pythonScript = "../vision/main.py"

func pythonBin() string {
	if bin := os.Getenv("PYTHON_BIN"); bin != "" {
		return bin
	}
	return "python3"
}

type Box struct {
	BBox      [4]int `json:"bbox"`
	IsChecked bool   `json:"is_checked"`
}

type DebugInfo struct {
	DetectionJobID string `json:"detection_job_id"`
}

type DetectResult struct {
	Boxes []Box      `json:"boxes"`
	Debug *DebugInfo `json:"debug,omitempty"`
}

func runDetect(file multipart.File, filename string, debug bool) (DetectResult, error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}
	tmp, err := os.CreateTemp("", "checkbox-*"+ext)
	if err != nil {
		return DetectResult{}, err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if _, err = io.Copy(tmp, file); err != nil {
		return DetectResult{}, err
	}
	tmp.Close()

	args := []string{pythonScript, tmp.Name()}

	var jobDir string
	if debug {
		if jobDir, err = os.MkdirTemp("", "checkbox-job-*"); err != nil {
			return DetectResult{}, err
		}
		args = append(args, jobDir)
	}

	cmd := exec.Command(pythonBin(), args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	var data []byte

	var runErr error
	if debug {
		runErr = cmd.Run()
	} else {
		data, runErr = cmd.Output()
	}
	if runErr != nil {
		fmt.Fprintf(os.Stderr, "detection engine error: %s\n", stderr.String())
		return DetectResult{}, fmt.Errorf("detection engine unavailable")
	}
	if debug {
		if data, err = os.ReadFile(filepath.Join(jobDir, "result.json")); err != nil {
			fmt.Fprintf(os.Stderr, "detection engine error: %s\n", err.Error())
			return DetectResult{}, fmt.Errorf("detection engine unavailable")
		}
	}

	var result DetectResult
	if err = json.Unmarshal(data, &result); err != nil {
		return DetectResult{}, err
	}
	if result.Boxes == nil {
		return DetectResult{}, fmt.Errorf("invalid detection output")
	}

	if debug {
		result.Debug = &DebugInfo{DetectionJobID: filepath.Base(jobDir)}
	}

	return result, nil
}
