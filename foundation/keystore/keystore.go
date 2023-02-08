// Package keystore implements the auth.Keystore interface. This implements
// an in-memory keystore for JWT support.
package keystore

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"
)

// Keystore represents an in memory store implementation of the
// KeyLookup interface for use with the auth package.
type Keystore struct {
	mu    sync.RWMutex
	store map[string]*rsa.PrivateKey
}

// New constructs an empty KeyStore ready for use.
func New() *Keystore {
	return &Keystore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

// NewMap constructs a KeyStore with an initial set of keys.
func NewMap(store map[string]*rsa.PrivateKey) *Keystore {
	return &Keystore{store: store}
}

// NewFS constructs a KeyStore based on a set of PEM files rooted inside
//
//	a directory. The name of each PEM file will be used as the key id.
func NewFS(fsys fs.FS) (*Keystore, error) {
	ks := Keystore{
		store: make(map[string]*rsa.PrivateKey),
	}

	fn := func(fileName string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsys.Open(fileName)
		if err != nil {
			return fmt.Errorf("opening key file: %w", err)
		}
		defer file.Close()

		privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("error reading pem file: %w", err)
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
		if err != nil {
			return fmt.Errorf("parsing auth private key: %w", err)
		}

		ks.store[strings.TrimSuffix(d.Name(), ".pem")] = privateKey

		ks.mu.Unlock()

		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return &ks, nil
}

// PrivateKey searches the key store for a given kid and returns
// the  private key
func (ks *Keystore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	ks.mu.Lock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, errors.New("kid lookup failed")
	}

	return privateKey, nil
}

// PublicKey searches the key store for a given id and returns
// the private key
func (ks *Keystore) PublicKey(kid string) (*rsa.PublicKey, error) {
	ks.mu.Lock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, errors.New("kid lookup failed")
	}

	return &privateKey.PublicKey, nil
}
