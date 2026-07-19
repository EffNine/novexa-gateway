package auth

import (
	"crypto/subtle"
	"errors"
)

// Service handles API key authentication
type Service struct {
	apiKey string
}

// NewService creates a new auth service
func NewService(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
	}
}

// Authenticate validates an API key
func (s *Service) Authenticate(providedKey string) error {
	if s.apiKey == "" {
		return errors.New("gateway API key not configured")
	}

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(s.apiKey), []byte(providedKey)) != 1 {
		return errors.New("invalid API key")
	}

	return nil
}

// IsConfigured returns true if an API key is set
func (s *Service) IsConfigured() bool {
	return s.apiKey != ""
}
