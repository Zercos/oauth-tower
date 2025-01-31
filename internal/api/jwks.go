package api

import (
	"crypto"
	"encoding/base64"
	"math/rand"
	"os"

	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt"
)

type KeyPair struct {
	Public  jose.JSONWebKey
	Private jose.JSONWebKey
}

type JWKManager struct {
	keyByKid map[string]KeyPair
}

func (m *JWKManager) LoadKeys() error {
	dirPath := os.Getenv("JWK_PATH")
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := file.Name()
		if file.IsDir() || (len(name) >= 10 && name[len(name)-10:] == "public.pem") {
			continue
		}
		privKeyFile := name
		pubKeyFile := name[:len(name)-4] + ".public.pem"
		pemPrivData, errPriv := os.ReadFile(dirPath + "/" + privKeyFile)
		pemPubData, errPub := os.ReadFile(dirPath + "/" + pubKeyFile)
		if errPriv != nil || errPub != nil {
			continue
		}
		pubKey, errPub := jwt.ParseRSAPublicKeyFromPEM(pemPubData)
		privKey, errPriv := jwt.ParseRSAPrivateKeyFromPEM(pemPrivData)
		if errPriv != nil || errPub != nil {
			continue
		}
		pubJWK := jose.JSONWebKey{Key: pubKey, Algorithm: "RS256"}
		privJWK := jose.JSONWebKey{Key: privKey, Algorithm: "RS256"}
		thumbprint, err := privJWK.Thumbprint(crypto.SHA256)
		if err != nil {
			continue
		}
		kid := base64.RawURLEncoding.EncodeToString(thumbprint[:])
		privJWK.KeyID = kid
		pubJWK.KeyID = kid
		m.keyByKid[kid] = KeyPair{Public: pubJWK, Private: privJWK}
	}
	return nil
}

func (m *JWKManager) GetSet() jose.JSONWebKeySet {
	keys := []jose.JSONWebKey{}
	for _, keyP := range m.keyByKid {
		keys = append(keys, keyP.Public)
	}
	return jose.JSONWebKeySet{Keys: keys}
}

func (m *JWKManager) GetSignKey() *jose.JSONWebKey {
	if len(m.keyByKid) == 0 {
		return nil
	}
	kids := []string{}
	for kid := range m.keyByKid {
		kids = append(kids, kid)
	}
	randInd := rand.Intn(len(kids))
	k := m.keyByKid[kids[randInd]]
	return &k.Private
}

func NewJWKManager() *JWKManager {
	return &JWKManager{keyByKid: make(map[string]KeyPair)}
}
