package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RemoteFromPath(path string) (string, error) {
	cmd := exec.Command("git", "-C", path, "config", "--get", "remote.origin.url")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(`executing "git -C %s config --get remote.origin.url": %v`, path, err)
	}
	return out.String(), nil
}

func RepoPathFromPath(path string) (string, error) {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--git-dir")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(`executing "git -C %s rev-parse --git-dir": %v`, path, err)
	}
	repoPath := strings.TrimSuffix(out.String(), ".git\n")
	if repoPath == "" {
		repoPath = path
	}
	return repoPath, nil
}
