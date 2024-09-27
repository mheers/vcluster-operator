package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// backend wraps the backend framework and adds a map for storing key value pairs
type backend struct {
	*framework.Backend
}

var _ logical.Factory = Factory

const configPath = "config"

// Factory configures and returns Mock backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := newBackend()
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func newBackend() (*backend, error) {
	b := &backend{}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(mockHelp),
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			b.paths(),
		),
	}

	return b, nil
}

func (b *backend) paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: framework.MatchAllRegex("path"),

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "Specifies the path of the secret.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRead,
					Summary:  "Retrieve the secret from the map.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDelete,
					Summary:  "Deletes the secret at the specified location.",
				},
			},

			ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

func (b *backend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}

type Config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (b *backend) readConfig(ctx context.Context, req *logical.Request) (*Config, error) {
	entry, err := req.Storage.Get(ctx, "/"+configPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result Config
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) writeConfig(ctx context.Context, req *logical.Request) (*Config, error) {
	config := &Config{
		URL:      req.Data["url"].(string),
		Username: req.Data["username"].(string),
		Password: req.Data["password"].(string),
	}
	// JSON encode the data
	buf, err := json.Marshal(config)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	// Store kv pairs in map at specified path
	entry := &logical.StorageEntry{
		Key:      "/" + configPath,
		Value:    buf,
		SealWrap: false,
	}
	if err = req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return config, nil
}

func (b *backend) handleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	config, err := b.readConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	path := data.Get("path").(string)

	// Decode the data
	rawData := make(map[string]interface{})
	rawData["path"] = path
	rawData["url"] = config.URL
	rawData["username"] = config.Username
	rawData["password"] = config.Password

	// Generate the response
	resp := &logical.Response{
		Data: rawData,
	}

	return resp, nil
}

func (b *backend) handleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	// Check to make sure that kv pairs provided
	if len(req.Data) == 0 {
		return nil, fmt.Errorf("data must be provided to store in secret")
	}

	path := data.Get("path").(string)

	if path != "config" {
		return nil, fmt.Errorf("path must be config")
	}

	config, err := b.writeConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", config)
	return nil, nil
}

func (b *backend) handleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	path := data.Get("path").(string)

	// Remove entry for specified path
	if err := req.Storage.Delete(ctx, req.ClientToken+"/"+path); err != nil {
		return nil, err
	}

	return nil, nil
}

const mockHelp = `
The Mock backend is a dummy secrets backend that stores kv pairs in a map.
`
