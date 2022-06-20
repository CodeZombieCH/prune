package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func TestPrune(t *testing.T) {
	repoPath := t.TempDir()

	createRepo(repoPath, t)

	// go run . -v -d 3 -m 2 -y 1 ~/temp/test-repo
	pruneArgs := []string{"-d", "3", "-m", "2", "-y", "1", repoPath}
	args := append([]string{"run", "./"}, pruneArgs...)
	cmd := exec.Command("go", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("prune failed with %v", err)
	}

	fmt.Print(string(out))

	filesToPrune := strings.Split(string(out), "\n")
	pruneFiles(filesToPrune, t)

	// Assert
	expected := []string{
		"2000-01-02T00-00-00.000Z",
		"2000-10-29T00-00-00.000Z",
		"2000-11-26T00-00-00.000Z",
		"2000-12-24T00-00-00.000Z",
		"2000-12-31T00-00-00.000Z",
		"2001-01-07T00-00-00.000Z",
	}

	files, err := os.ReadDir(repoPath)
	if err != nil {
		t.Errorf("Failed to read files in repo: %v", err)
	}

	if expected, actual := 6, len(files); actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	for _, v := range expected {
		if _, err := os.Stat(path.Join(repoPath, v)); err != nil && errors.Is(err, os.ErrNotExist) {
			t.Errorf("Expected file %v does not exist", v)
		}
	}
}

func createRepo(repoPath string, t *testing.T) {
	args := []string{repoPath}
	goArgs := append([]string{"run", "./cmd/test-repo"}, args...)

	cmd := exec.Command("go", goArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("prune failed with %v", err)
	}

	fmt.Print(string(out))
}

func pruneFiles(filesToPrune []string, t *testing.T) {
	for _, file := range filesToPrune {
		cmd := exec.Command("rm", "-rf", file)
		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("rm -rf failed with %v", err)
		}
	}
}
