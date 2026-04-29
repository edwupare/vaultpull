package sync

import (
	"fmt"
	"os"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/envwriter"
	"github.com/user/vaultpull/internal/vault"
)

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	cfg    *config.Config
	client *vault.Client
}

// New creates a new Syncer from the given config.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}
	return &Syncer{cfg: cfg, client: client}, nil
}

// Run performs the full sync: read secrets, diff, backup, merge, and write.
func (s *Syncer) Run() error {
	secrets, err := s.client.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("reading secrets: %w", err)
	}

	filtered := envwriter.FilterSecrets(secrets, s.cfg.Namespace)

	// Load existing env for diffing
	existing := map[string]string{}
	if data, err := envwriter.ReadEnvFile(s.cfg.OutputFile); err == nil {
		existing = data
	}

	diff := envwriter.DiffEnv(existing, filtered)
	if !diff.HasChanges() {
		fmt.Fprintln(os.Stdout, "No changes detected.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Changes:\n%s", diff.Summary())

	if err := envwriter.BackupFile(s.cfg.OutputFile, s.cfg.BackupCount); err != nil {
		return fmt.Errorf("backing up env file: %w", err)
	}

	if err := envwriter.MergeEnvFile(s.cfg.OutputFile, filtered); err != nil {
		return fmt.Errorf("merging env file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Synced %d secret(s) to %s\n", len(filtered), s.cfg.OutputFile)
	return nil
}
