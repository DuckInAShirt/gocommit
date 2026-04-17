package git

import (
	"fmt"
	"os/exec"
	"strings"
)

var ExcludePatterns = []string{
	"*.lock",
	"*.lockb",
	"*-lock.json",
	"*-lock.yaml",
	"package-lock.json",
	"pnpm-lock.yaml",
	"yarn.lock",
	"go.sum",
	"*.min.js",
	"*.min.css",
}

func GetStagedDiff() (string, error) {
	args := []string{"diff", "--cached"}
	for _, p := range ExcludePatterns {
		args = append(args, fmt.Sprintf(":%s", p))
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git diff failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("git diff failed: %w", err)
	}

	diff := strings.TrimSpace(string(out))
	if diff == "" {
		return "", fmt.Errorf("no staged changes found, run `git add` first")
	}

	return diff, nil
}

func GetDiffStat() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--stat")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff --stat failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

func HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	return cmd.Run() != nil
}

func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func AmendCommit(message string) error {
	cmd := exec.Command("git", "commit", "--amend", "-m", message)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit --amend failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func StageAll() error {
	cmd := exec.Command("git", "add", "-A")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git add -A failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func TruncateDiff(diff string, maxChars int) string {
	if len(diff) <= maxChars {
		return diff
	}
	truncated := diff[:maxChars]
	if idx := strings.LastIndex(truncated, "\n"); idx > 0 {
		truncated = truncated[:idx]
	}
	return truncated + "\n... (diff truncated due to size)"
}
