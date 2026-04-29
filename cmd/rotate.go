package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/sync"
)

var (
	rotateDryRun    bool
	rotateNamespace string
	rotateMaxBackup int
	rotateAudit     string
)

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate secrets: backup existing .env, fetch from Vault, merge and write",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		syncer, err := sync.New(cfg)
		if err != nil {
			return fmt.Errorf("syncer: %w", err)
		}

		rotator := sync.NewRotator(syncer, sync.RotatorConfig{
			VaultPath:  cfg.SecretPath,
			OutputPath: cfg.OutputPath,
			Namespace:  rotateNamespace,
			MaxBackups: rotateMaxBackup,
			DryRun:     rotateDryRun,
			AuditPath:  rotateAudit,
		})

		result, err := rotator.Run()
		if err != nil {
			return fmt.Errorf("rotate: %w", err)
		}

		if rotateDryRun {
			fmt.Fprintln(os.Stdout, "[dry-run] no changes written")
		} else {
			fmt.Fprintf(os.Stdout, "Rotated %s\n", cfg.OutputPath)
			if result.BackupPath != "" {
				fmt.Fprintf(os.Stdout, "Backup:  %s\n", result.BackupPath)
			}
		}

		fmt.Fprintln(os.Stdout, result.Changes.Summary())
		return nil
	},
}

func init() {
	rotateCmd.Flags().BoolVar(&rotateDryRun, "dry-run", false, "Preview changes without writing")
	rotateCmd.Flags().StringVar(&rotateNamespace, "namespace", "", "Filter secrets by key namespace prefix")
	rotateCmd.Flags().IntVar(&rotateMaxBackup, "max-backups", 5, "Maximum number of backup files to keep")
	rotateCmd.Flags().StringVar(&rotateAudit, "audit", "", "Path to audit log file (stdout if empty)")
	rootCmd.AddCommand(rotateCmd)
}
