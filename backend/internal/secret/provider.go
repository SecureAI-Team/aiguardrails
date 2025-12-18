package secret

import (
	"context"
	"errors"
	"os"
)

// Provider defines how to retrieve secrets by key.
type Provider interface {
	GetSecret(ctx context.Context, key string) (string, error)
}

// EnvProvider reads secrets from environment variables.
type EnvProvider struct{}

func (EnvProvider) GetSecret(_ context.Context, key string) (string, error) {
	if key == "" {
		return "", errors.New("secret key empty")
	}
	v := os.Getenv(key)
	if v == "" {
		return "", errors.New("secret not found")
	}
	return v, nil
}

// StaticProvider returns a fixed map (useful for tests).
type StaticProvider struct {
	Values map[string]string
}

func (s StaticProvider) GetSecret(_ context.Context, key string) (string, error) {
	if v, ok := s.Values[key]; ok && v != "" {
		return v, nil
	}
	return "", errors.New("secret not found")
}

