package repo

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// checks if the given path is a git url
func IsUpsteam(path string) bool {
	if strings.HasPrefix(path, "git@") {
		return true
	}
	if strings.HasPrefix(path, "https://") {
		return true
	}
	return false
}

// clone a git repo to the /tmp/<repo>, striping the .git suffix
// this has dependencies on the git binary
func Clone(ctx context.Context, repo string) (path string, err error) {
	// get the repo path
	path, err = constructPath(repo)
	if err != nil {
		return "", err
	}

	// create repo dir if it doesn't exist
	os.MkdirAll(path, 0755)

	fmt.Printf("Cloning repo %s to path %s\n", repo, path)

	// clone the repo
	cmd := exec.CommandContext(ctx, "git", "clone", "--recursive", repo, path)
	err = cmd.Run()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err != nil {
		return "", fmt.Errorf("error cloning repo: %s", err)
	}

	return path, err
}

// construct the path to the repo using the os.TempDir and the repo name
func constructPath(path string) (repoPath string, err error) {
	repoWithoutDotGit := strings.TrimSuffix(path, ".git")
	repoPath = filepath.Join(os.TempDir(), filepath.Base(repoWithoutDotGit))
	return repoPath, nil
}
