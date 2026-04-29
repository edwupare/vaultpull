package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/envwriter"
)

// RotatorConfig holds configuration for the rotation workflow.
type RotatorConfig struct {
	VaultPath   string
	OutputPath  string
	Namespace   string
	MaxBackups  int
	DryRun      bool
	AuditPath   string
}

// Rotator orchestrates fetching secrets from Vault and rotating the env file.
type Rotator struct {
	cfg    RotatorConfig
	syncer *Syncer
}

// NewRotator creates a Rotator using the provided Syncer and config.
func NewRotator(s *Syncer, cfg RotatorConfig) *Rotator {
	return &Rotator{syncer: s, cfg: cfg}
}

// Run fetches secrets and rotates the target env file.
func (r *Rotator) Run() (*envwriter.RotateResult, error) {
	secrets, err := r.syncer.client.ReadSecrets(r.cfg.VaultPath)
	if err != nil {
		return nil, fmt.Errorf("rotator: failed to read secrets: %w", err)
	}

	var auditLogger *envwriter.AuditLogger
	if r.cfg.AuditPath != "" {
		auditLogger, err = envwriter.NewAuditLogger(r.cfg.AuditPath)
		if err != nil {
			return nil, fmt.Errorf("rotator: failed to init audit logger: %w", err)
		}
	}

	opts := envwriter.RotateOptions{
		Namespace:   r.cfg.Namespace,
		MaxBackups:  r.cfg.MaxBackups,
		AuditLogger: auditLogger,
		DryRun:      r.cfg.DryRun,
	}

	result, err := envwriter.RotateEnvFile(r.cfg.OutputPath, secrets, opts)
	if err != nil {
		return nil, fmt.Errorf("rotator: rotation failed: %w", err)
	}

	return result, nil
}
