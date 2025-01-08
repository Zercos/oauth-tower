package api

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt"
)

type JWKManager struct {
	keyByKid map[string]jose.JSONWebKey
}

func (m *JWKManager) ImportKeys() error {
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
				return err
			}
			pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pemData)
			if err != nil {
				return err
			}
			keyIdHash := md5.Sum(pemData)
			keyId := hex.EncodeToString(keyIdHash[:])
			m.keyByKid[keyId] = jose.JSONWebKey{Key: pubKey, KeyID: keyId, Algorithm: "RS256"}
		}
	}
	return nil
}
