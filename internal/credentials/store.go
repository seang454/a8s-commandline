package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/zalando/go-keyring"
)

const serviceName = "a8s-cli"

var ErrNotFound = errors.New("credential not found")

type Record struct {
	AccessToken        string    `json:"accessToken"`
	RefreshToken       string    `json:"refreshToken,omitempty"`
	IDToken            string    `json:"idToken,omitempty"`
	AccessTokenExpiry  time.Time `json:"accessTokenExpiry"`
	RefreshTokenExpiry time.Time `json:"refreshTokenExpiry,omitempty"`
	Issuer             string    `json:"issuer"`
	ClientID           string    `json:"clientId"`
	Subject            string    `json:"subject,omitempty"`
	Username           string    `json:"username,omitempty"`
	Email              string    `json:"email,omitempty"`
	Roles              []string  `json:"roles,omitempty"`
}

type Store interface {
	Get(key string) (Record, error)
	Set(key string, record Record) error
	Delete(key string) error
}

// NativeStore uses the operating-system keyring and falls back to a restricted
// local file when a keyring service is unavailable.
type NativeStore struct {
	fallbackPath string
	warningOut   io.Writer
	warningOnce  sync.Once
}

func NewNativeStore() (*NativeStore, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve credential directory: %w", err)
	}
	return &NativeStore{fallbackPath: filepath.Join(dir, "a8s", "credentials.json")}, nil
}

func (s *NativeStore) WarnFallbackTo(writer io.Writer) {
	s.warningOut = writer
}

func (s *NativeStore) Get(key string) (Record, error) {
	value, err := keyring.Get(serviceName, key)
	if err == nil {
		return decode(value)
	}
	if !errors.Is(err, keyring.ErrNotFound) {
		record, fallbackErr := s.getFallback(key)
		if fallbackErr == nil {
			s.warnFallback()
		}
		return record, fallbackErr
	}
	record, fallbackErr := s.getFallback(key)
	if fallbackErr == nil {
		s.warnFallback()
		return record, nil
	}
	if errors.Is(fallbackErr, ErrNotFound) {
		return Record{}, ErrNotFound
	}
	return Record{}, fallbackErr
}

func (s *NativeStore) Set(key string, record Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode credentials: %w", err)
	}
	if err := keyring.Set(serviceName, key, string(data)); err == nil {
		_ = s.deleteFallback(key)
		return nil
	}
	s.warnFallback()
	return s.setFallback(key, record)
}

func (s *NativeStore) Delete(key string) error {
	keyringErr := keyring.Delete(serviceName, key)
	fileErr := s.deleteFallback(key)
	if keyringErr == nil || errors.Is(keyringErr, keyring.ErrNotFound) {
		return fileErr
	}
	if fileErr == nil {
		return nil
	}
	return fmt.Errorf("delete credentials: keyring: %v; fallback: %v", keyringErr, fileErr)
}

func (s *NativeStore) getFallback(key string) (Record, error) {
	records, err := s.loadFallback()
	if err != nil {
		return Record{}, err
	}
	record, ok := records[key]
	if !ok {
		return Record{}, ErrNotFound
	}
	return record, nil
}

func (s *NativeStore) setFallback(key string, record Record) error {
	records, err := s.loadFallback()
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	records[key] = record
	return s.saveFallback(records)
}

func (s *NativeStore) deleteFallback(key string) error {
	records, err := s.loadFallback()
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	delete(records, key)
	return s.saveFallback(records)
}

func (s *NativeStore) loadFallback() (map[string]Record, error) {
	data, err := os.ReadFile(s.fallbackPath)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]Record{}, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read credential fallback: %w", err)
	}
	var records map[string]Record
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("decode credential fallback: %w", err)
	}
	return records, nil
}

func (s *NativeStore) saveFallback(records map[string]Record) error {
	if err := os.MkdirAll(filepath.Dir(s.fallbackPath), 0o700); err != nil {
		return fmt.Errorf("create credential fallback directory: %w", err)
	}
	data, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("encode credential fallback: %w", err)
	}
	temporary := s.fallbackPath + ".tmp"
	if err := os.WriteFile(temporary, data, 0o600); err != nil {
		return fmt.Errorf("write credential fallback: %w", err)
	}
	if err := os.Rename(temporary, s.fallbackPath); err != nil {
		return fmt.Errorf("replace credential fallback: %w", err)
	}
	return nil
}

func (s *NativeStore) warnFallback() {
	if s.warningOut == nil {
		return
	}
	s.warningOnce.Do(func() {
		fmt.Fprintf(s.warningOut, "Warning: operating-system credential storage is unavailable; using restricted file %s\n", s.fallbackPath)
	})
}

func decode(value string) (Record, error) {
	var record Record
	if err := json.Unmarshal([]byte(value), &record); err != nil {
		return Record{}, fmt.Errorf("decode credentials: %w", err)
	}
	return record, nil
}

func Key(contextName, configured string) string {
	if strings.TrimSpace(configured) != "" {
		return strings.TrimSpace(configured)
	}
	return "context:" + contextName
}
