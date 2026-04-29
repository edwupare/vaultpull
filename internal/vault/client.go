package vault

import (
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with project-specific helpers.
type Client struct {
	api    *vaultapi.Client
	prefix string
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(addr, token, prefix string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to create client: %w", err)
	}

	api.SetToken(token)

	return &Client{
		api:    api,
		prefix: strings.TrimRight(prefix, "/"),
	}, nil
}

// ReadSecrets reads all key-value pairs from the given secret path.
// The path is relative to the configured prefix.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	fullPath := c.prefix + "/" + strings.TrimLeft(path, "/")

	secret, err := c.api.Logical().Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to read path %q: %w", fullPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: no secret found at path %q", fullPath)
	}

	data, ok := secret.Data["data"]
	if !ok {
		// KV v1 — data is at the top level
		return flattenData(secret.Data), nil
	}

	// KV v2 — data is nested under "data"
	kvData, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault: unexpected data format at path %q", fullPath)
	}
	return flattenData(kvData), nil
}

func flattenData(raw map[string]interface{}) map[string]string {
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}
