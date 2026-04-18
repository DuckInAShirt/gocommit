package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xinranzhao/gocommit/internal/commit"
	"github.com/xinranzhao/gocommit/internal/config"
)

var (
	autoStage bool
	amend     bool
	autoYes   bool
	dryRun    bool
	debug     bool
)

var rootCmd = &cobra.Command{
	Use:   "gocommit",
	Short: "AI-powered Chinese git commit message generator",
	Long:  "gocommit analyzes your staged changes and generates Chinese conventional commit messages using AI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commit.Run(commit.Options{
			AutoStage: autoStage,
			Amend:     amend,
			AutoYes:   autoYes,
			DryRun:    dryRun,
			Debug:     debug,
		})
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure gocommit settings",
	Long:  "Set or view gocommit configuration (API key, model, base URL).",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			masked := maskKey(cfg.APIKey)
			fmt.Printf("API Key:  %s\n", masked)
			fmt.Printf("Base URL: %s\n", cfg.BaseURL)
			fmt.Printf("Model:    %s\n", cfg.Model)
			return nil
		}

		for _, arg := range args {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid config format, use key=value (e.g. api_key=sk-xxx)")
			}
			key, value := parts[0], parts[1]
			switch key {
			case "api_key":
				cfg.APIKey = value
			case "base_url":
				cfg.BaseURL = value
			case "model":
				cfg.Model = value
			default:
				return fmt.Errorf("unknown config key: %s (valid: api_key, base_url, model)", key)
			}
		}

		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println("Config saved.")
		return nil
	},
}

func maskKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func init() {
	rootCmd.Flags().BoolVarP(&autoStage, "all", "a", false, "stage all changes before generating")
	rootCmd.Flags().BoolVarP(&amend, "amend", "", false, "amend the previous commit")
	rootCmd.Flags().BoolVarP(&autoYes, "yes", "y", false, "skip confirmation and auto-commit")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "show message without committing")
	rootCmd.Flags().BoolVarP(&debug, "debug", "", false, "print API request/response for debugging")

	rootCmd.AddCommand(configCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
