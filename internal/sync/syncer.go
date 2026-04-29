package sync

import (
	"fmt"
	"os"

	"github.com/vaultpull/internal/config"
	"github.com/vaultpull/internal/envwriter"
	"github.com/vaultpull/internal/vault"
)

// Result holds the outcome of a sync operation.
type Result struct {
	Path      string
	Written   int
	Filtered  int
	OutputFile string
}

// Syncer orchestrates reading secrets from Vault and writing them to a file.
type Syncer struct {
	client *vault.Client
	cfg    *config.Config
}

// New creates a new Syncer from the provided config.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create vault client: %w", err)
	}
	return &Syncer{client: client, cfg: cfg}, nil
}

// Run performs the full sync: read secrets, filter by namespace, write to output.
func (s *Syncer) Run(secretPath, namespace, outputFile string) (*Result, error) {
	secrets, err := s.client.ReadSecrets(secretPath)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to read secrets at %q: %w", secretPath, err)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create output file %q: %w", outputFile, err)
	}
	defer f.Close()

	written, filtered, err := envwriter.FilteredWrite(f, secrets, namespace)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to write env file: %w", err)
	}

	return &Result{
		Path:       secretPath,
		Written:    written,
		Filtered:   filtered,
		OutputFile: outputFile,
	}, nil
}
