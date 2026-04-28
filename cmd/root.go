package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/config"
)

var (
	outputFile string
	namespace  string
	prefix     string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull fetches secrets from a HashiCorp Vault KV store
and writes them to a local .env file, with optional namespace
filtering and key prefix support.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputFile != "" {
			os.Setenv("VAULTPULL_OUTPUT", outputFile)
		}
		if namespace != "" {
			os.Setenv("VAULT_NAMESPACE", namespace)
		}
		if prefix != "" {
			os.Setenv("VAULTPULL_PREFIX", prefix)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration error: %w", err)
		}

		fmt.Printf("vaultpull starting\n")
		fmt.Printf("  Vault address : %s\n", cfg.VaultAddr)
		fmt.Printf("  Namespace     : %s\n", cfg.Namespace)
		fmt.Printf("  Output file   : %s\n", cfg.OutputFile)
		fmt.Printf("  Key prefix    : %s\n", cfg.Prefix)

		// Sync logic will be wired here in subsequent phases.
		fmt.Println("(sync not yet implemented)")
		return nil
	},
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output .env file path (overrides VAULTPULL_OUTPUT)")
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Vault namespace to filter secrets (overrides VAULT_NAMESPACE)")
	rootCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "prefix to add to every exported key (overrides VAULTPULL_PREFIX)")
}
