package commit

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xinranzhao/gocommit/internal/ai"
	"github.com/xinranzhao/gocommit/internal/config"
	"github.com/xinranzhao/gocommit/internal/git"
)

type Options struct {
	AutoStage bool
	Amend     bool
	AutoYes   bool
	DryRun    bool
	Debug     bool
}

func Run(opts Options) error {
	if !git.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	if opts.AutoStage {
		if err := git.StageAll(); err != nil {
			return err
		}
	}

	diff, err := git.GetStagedDiff()
	if err != nil {
		return err
	}

	stat, _ := git.GetDiffStat()
	if stat != "" {
		fmt.Println(stat)
		fmt.Println()
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	if cfg.APIKey == "" {
		return fmt.Errorf("API key not set, run `gocommit config` or set OPENAI_API_KEY env")
	}

	client := ai.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model, opts.Debug)

	fmt.Print("Generating commit message... ")
	message, err := client.GenerateCommitMessage(diff)
	if err != nil {
		return fmt.Errorf("generate commit message failed: %w", err)
	}
	fmt.Println("done!")
	fmt.Println()

	fmt.Printf("  %s\n\n", message)

	if opts.DryRun {
		return nil
	}

	if opts.AutoYes {
		return doCommit(message, opts.Amend)
	}

	fmt.Print("Commit with this message? [y/e/r/n] (y=yes, e=edit, r=retry, n=abort): ")
	action, err := readAction()
	if err != nil {
		return err
	}

	switch action {
	case "y", "":
		return doCommit(message, opts.Amend)
	case "e":
		edited := editMessage(message)
		if edited == "" {
			fmt.Println("Aborted.")
			return nil
		}
		return doCommit(edited, opts.Amend)
	case "r":
		fmt.Println("Regenerating...")
		return Run(opts)
	default:
		fmt.Println("Aborted.")
		return nil
	}
}

func doCommit(message string, amend bool) error {
	if amend {
		return git.AmendCommit(message)
	}
	return git.Commit(message)
}

func readAction() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.ToLower(input)), nil
}

func editMessage(original string) string {
	fmt.Printf("Edit message (press Enter to keep, empty to abort):\n  ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return original
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	return input
}
