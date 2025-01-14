package api

import (
	"crypto"
	"encoding/base64"
	"os"

	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt"
)

type JWKManager struct {
	keyByKid map[string]jose.JSONWebKey
}

func (m *JWKManager) LoadKeys() error {
	dirPath := "./keys"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := file.Name()
		if !file.IsDir() && len(name) >= 10 && name[len(name)-10:] == "public.pem" {
			pemData, err := os.ReadFile(dirPath + "/" + name)
			if err != nil {
				continue
			}
			pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pemData)
			if err != nil {
				continue
			}
			key := jose.JSONWebKey{Key: pubKey, Algorithm: "RS256"}
			thumbprint, err := key.Thumbprint(crypto.SHA256)
			if err != nil {
				continue
			}
			kid := base64.RawURLEncoding.EncodeToString(thumbprint[:])
			key.KeyID = kid
			m.keyByKid[kid] = key
		}
	}
	return nil
}

func (m *JWKManager) GetSet() jose.JSONWebKeySet {
	keys := []jose.JSONWebKey{}
	for _, key := range m.keyByKid {
		keys = append(keys, key)
	}
	return jose.JSONWebKeySet{Keys: keys}
}

func (m *JWKManager) GetSignKey() jose.JSONWebKey {
	var key jose.JSONWebKey
	for _, key := range m.keyByKid {
		return key
	}
	return key
}

func NewJWKManager() *JWKManager {
	return &JWKManager{keyByKid: make(map[string]jose.JSONWebKey)}
}
