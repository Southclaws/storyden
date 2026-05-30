package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-yaml"
)

const (
	configDirectory = "storyden"
	configFilename  = "config.yaml"
)

type Config struct {
	CurrentContext string             `yaml:"current_context,omitempty"`
	Contexts       map[string]Context `yaml:"contexts,omitempty"`
}

type AuthStorage string
type AuthMethod string

const (
	AuthStorageCredentialStore AuthStorage = "credential_store"
	AuthStorageFile            AuthStorage = "file"

	AuthMethodOAuthDevice AuthMethod = "oauth_device"
	AuthMethodAccessKey   AuthMethod = "access_key"
)

type Context struct {
	APIURL   string      `yaml:"api_url"`
	AuthType AuthStorage `yaml:"auth_type,omitempty"`
	Auth     *Auth       `yaml:"auth,omitempty"`
}

type Auth struct {
	Method       AuthMethod `yaml:"method,omitempty"`
	AccessToken  string     `yaml:"access_token,omitempty"`
	RefreshToken string     `yaml:"refresh_token,omitempty"`
	TokenType    string     `yaml:"token_type,omitempty"`
	ExpiresAt    time.Time  `yaml:"expires_at,omitempty"`
	Scope        string     `yaml:"scope,omitempty"`
	Issuer       string     `yaml:"issuer,omitempty"`
	ClientID     string     `yaml:"client_id,omitempty"`
}

func (a Auth) MethodOrDefault() AuthMethod {
	if a.Method != "" {
		return a.Method
	}

	return AuthMethodOAuthDevice
}

type Store struct {
	path       string
	credential CredentialStore
}

func NewStore() (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	return NewStoreAt(filepath.Join(configDir, configDirectory, configFilename)), nil
}

func NewStoreAt(path string) *Store {
	return NewStoreAtWithCredentialStore(path, newKeyringStore())
}

func NewFileStoreAt(path string) *Store {
	return NewStoreAtWithCredentialStore(path, unavailableCredentialStore{})
}

func NewStoreAtWithCredentialStore(path string, credential CredentialStore) *Store {
	if credential == nil {
		credential = unavailableCredentialStore{}
	}
	return &Store{
		path:       path,
		credential: credential,
	}
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) CredentialStoreAvailable() bool {
	return s.credential.Available()
}

func (s *Store) DeleteAuth(contextName string) error {
	return s.credential.DeleteAuth(contextName)
}

func (s *Store) DefaultAuthStorage() AuthStorage {
	if s.CredentialStoreAvailable() {
		return AuthStorageCredentialStore
	}

	return AuthStorageFile
}

func (s *Store) Load() (*Config, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(), nil
		}

		return nil, err
	}

	if len(data) == 0 {
		return New(), nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.normalise()
	if err := s.loadCredentials(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (s *Store) Save(cfg *Config) error {
	cfg.normalise()

	if err := s.saveCredentials(cfg); err != nil {
		return err
	}

	persisted := cfg.persisted()
	data, err := yaml.Marshal(persisted)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}

	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return err
	}

	return nil
}

func New() *Config {
	return &Config{
		Contexts: map[string]Context{},
	}
}

func (c *Config) UpsertContext(name string, context Context) {
	c.normalise()
	c.Contexts[name] = context
}

func (c *Config) SetCurrentContext(name string) {
	c.CurrentContext = name
}

func (c *Config) normalise() {
	if c.Contexts == nil {
		c.Contexts = map[string]Context{}
	}

	for name, ctx := range c.Contexts {
		if ctx.AuthType == "" && ctx.Auth != nil {
			ctx.AuthType = AuthStorageFile
			c.Contexts[name] = ctx
		}
	}
}

func (c *Config) persisted() *Config {
	result := &Config{
		CurrentContext: c.CurrentContext,
		Contexts:       map[string]Context{},
	}

	for name, ctx := range c.Contexts {
		if ctx.AuthType == "" && ctx.Auth != nil {
			ctx.AuthType = AuthStorageFile
		}
		if ctx.AuthType == AuthStorageCredentialStore {
			ctx.Auth = nil
		}

		result.Contexts[name] = ctx
	}

	return result
}

func (s *Store) saveCredentials(cfg *Config) error {
	for name, ctx := range cfg.Contexts {
		switch ctx.AuthType {
		case AuthStorageCredentialStore:
			if !s.CredentialStoreAvailable() {
				return fmt.Errorf("credential store is not available for context %q", name)
			}
			if ctx.Auth == nil {
				if err := s.credential.DeleteAuth(name); err != nil {
					return err
				}
				continue
			}
			if err := s.credential.SetAuth(name, *ctx.Auth); err != nil {
				return err
			}

		case AuthStorageFile, "":
			if err := s.credential.DeleteAuth(name); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported auth_type %q for context %q", ctx.AuthType, name)
		}
	}

	return nil
}

func (s *Store) loadCredentials(cfg *Config) error {
	for name, ctx := range cfg.Contexts {
		if ctx.AuthType != AuthStorageCredentialStore {
			continue
		}

		ctx.Auth = nil
		auth, ok, err := s.credential.GetAuth(name)
		if err != nil {
			return err
		}
		if ok {
			ctx.Auth = &auth
		}

		cfg.Contexts[name] = ctx
	}

	return nil
}
