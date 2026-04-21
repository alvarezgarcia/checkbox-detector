package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

const pythonBin = "../vision/venv/bin/python"
const pythonScript = "../vision/main.py"

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

func runDetect(file multipart.File, filename string, debug bool) (result DetectResult, err error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}
	tmp, err := os.CreateTemp("", "checkbox-*"+ext)
	if err != nil {
		return
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if _, err = io.Copy(tmp, file); err != nil {
		return
	}
	tmp.Close()

	args := []string{pythonScript, tmp.Name()}

	var jobDir string
	if debug {
		if jobDir, err = os.MkdirTemp("", "checkbox-job-*"); err != nil {
			return
		}
		args = append(args, jobDir)
	}

	cmd := exec.Command(pythonBin, args...)
	cmd.Stderr = os.Stderr

	var data []byte

	if debug {
		if err = cmd.Run(); err != nil {
			return
		}
		data, err = os.ReadFile(filepath.Join(jobDir, "result.json"))
	} else {
		data, err = cmd.Output()
	}
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return
	}
	if result.Boxes == nil {
		err = fmt.Errorf("invalid detection output")
		return
	}

	if debug {
		result.Debug = &DebugInfo{DetectionJobID: filepath.Base(jobDir)}
	}

	return
}
