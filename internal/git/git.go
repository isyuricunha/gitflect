package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Mirror performs a mirror clone from cloneURL to a temporary directory,
// attempts an LFS fetch, pushes the mirror to pushURL, and pushes LFS to pushURL.
func Mirror(cloneURL, pushURL, repoName string, secrets []string) error {
	workDir := filepath.Join(os.TempDir(), "gitflect_"+repoName)
	// Clean the directory to avoid leftover files
	os.RemoveAll(workDir)
	defer os.RemoveAll(workDir)

	// 1. Clone --mirror
	if err := runRedacted(secrets, "", "git", "clone", "--mirror", cloneURL, workDir); err != nil {
		return fmt.Errorf("git clone error: %w", err)
	}

	// 2. Attempt a git lfs fetch --all to ensure LFS objects are downloaded.
	// Executed silently because if the repo does not use LFS, it shouldn't be a fatal error.
	_ = runRedacted(secrets, workDir, "git", "lfs", "fetch", "--all")

	// 3. Push LFS first!
	// We must push LFS objects before pushing the mirror itself,
	// because GitLab will reject the git push if commits reference missing LFS blocks.
	_ = runRedacted(secrets, workDir, "git", "lfs", "push", "--all", pushURL)

	// 4. Push --mirror to destination (commits and refs)
	if err := runRedacted(secrets, workDir, "git", "push", "--mirror", pushURL); err != nil {
		return fmt.Errorf("git push mirror error: %w", err)
	}

	return nil
}

// runRedacted executes a command capturing stdout/stderr and replacing any secrets with ***
func runRedacted(secrets []string, dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	if err != nil {
		outputStr := outBuf.String()
		for _, secret := range secrets {
			if secret == "" {
				continue
			}
			outputStr = strings.ReplaceAll(outputStr, secret, "***")
		}
		
		// Also hide secrets in arguments in case the error exposes the entire command string
		errStr := err.Error()
		for _, secret := range secrets {
			if secret == "" {
				continue
			}
			errStr = strings.ReplaceAll(errStr, secret, "***")
		}
		return fmt.Errorf("exec error: %s - output: %s", errStr, outputStr)
	}
	return nil
}
