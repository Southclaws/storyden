package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"

	"github.com/99designs/keyring"
)

const keyringServiceName = "storyden"

type CredentialStore interface {
	SetAuth(contextName string, auth Auth) error
	GetAuth(contextName string) (Auth, bool, error)
	DeleteAuth(contextName string) error
	Available() bool
}

type unavailableCredentialStore struct{}

func (unavailableCredentialStore) SetAuth(_ string, _ Auth) error {
	return fmt.Errorf("credential store is not available")
}

func (unavailableCredentialStore) GetAuth(_ string) (Auth, bool, error) {
	return Auth{}, false, nil
}

func (unavailableCredentialStore) DeleteAuth(_ string) error { return nil }

func (unavailableCredentialStore) Available() bool { return false }

type keyringCredentialStore struct {
	keyring keyring.Keyring
}

func newKeyringStore() CredentialStore {
	backends := nativeKeyringBackends()
	if len(backends) == 0 {
		return unavailableCredentialStore{}
	}

	kr, err := keyring.Open(keyring.Config{
		ServiceName:     keyringServiceName,
		AllowedBackends: backends,
	})
	if err != nil {
		return unavailableCredentialStore{}
	}

	return keyringCredentialStore{keyring: kr}
}

func nativeKeyringBackends() []keyring.BackendType {
	switch runtime.GOOS {
	case "darwin":
		return []keyring.BackendType{keyring.KeychainBackend}
	case "windows":
		return []keyring.BackendType{keyring.WinCredBackend}
	case "linux":
		return []keyring.BackendType{
			keyring.SecretServiceBackend,
			keyring.KWalletBackend,
			keyring.PassBackend,
			keyring.KeyCtlBackend,
		}
	default:
		return nil
	}
}

func (k keyringCredentialStore) SetAuth(contextName string, auth Auth) error {
	data, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	if err := k.keyring.Set(keyring.Item{Key: authKey(contextName), Data: data}); err != nil {
		return fmt.Errorf("store auth credentials for context %q: %w", contextName, err)
	}

	return nil
}

func (k keyringCredentialStore) GetAuth(contextName string) (Auth, bool, error) {
	item, err := k.keyring.Get(authKey(contextName))
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return Auth{}, false, nil
		}

		return Auth{}, false, fmt.Errorf("load auth credentials for context %q: %w", contextName, err)
	}

	var auth Auth
	if err := json.Unmarshal(item.Data, &auth); err != nil {
		return Auth{}, false, fmt.Errorf("decode auth credentials for context %q: %w", contextName, err)
	}

	return auth, true, nil
}

func (k keyringCredentialStore) DeleteAuth(contextName string) error {
	err := k.keyring.Remove(authKey(contextName))
	if errors.Is(err, keyring.ErrKeyNotFound) || err == nil {
		return nil
	}

	return fmt.Errorf("delete auth credentials for context %q: %w", contextName, err)
}

func authKey(contextName string) string {
	return "context:" + contextName
}

func (k keyringCredentialStore) Available() bool { return true }
