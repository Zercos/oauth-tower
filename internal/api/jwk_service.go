package api

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"sync"

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
	dirPath := config.getJwkPath()
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	var keysMap sync.Map
	var wg sync.WaitGroup

	for _, file := range files {
		name := file.Name()
		if file.IsDir() || (len(name) >= 10 && name[len(name)-10:] == "public.pem") {
			continue
		}
		wg.Add(1)
		go func(name string, keys *sync.Map, wg *sync.WaitGroup) {
			defer wg.Done()
			privKeyFile := name
			pubKeyFile := name[:len(name)-4] + ".public.pem"
			pemPrivData, errPriv := os.ReadFile(dirPath + "/" + privKeyFile)
			pemPubData, errPub := os.ReadFile(dirPath + "/" + pubKeyFile)
			if errPriv != nil || errPub != nil {
				return
			}
			pubKey, errPub := jwt.ParseRSAPublicKeyFromPEM(pemPubData)
			privKey, errPriv := jwt.ParseRSAPrivateKeyFromPEM(pemPrivData)
			if errPriv != nil || errPub != nil {
				return
			}
			pubJWK := jose.JSONWebKey{Key: pubKey, Algorithm: "RS256"}
			privJWK := jose.JSONWebKey{Key: privKey, Algorithm: "RS256"}
			thumbprint, err := privJWK.Thumbprint(crypto.SHA256)
			if err != nil {
				return
			}
			kid := base64.RawURLEncoding.EncodeToString(thumbprint[:])
			privJWK.KeyID = kid
			pubJWK.KeyID = kid
			keys.Store(kid, KeyPair{Public: pubJWK, Private: privJWK})
		}(name, &keysMap, &wg)

	}

	wg.Wait()
	keysMap.Range(func(key, value interface{}) bool {
		m.keyByKid[key.(string)] = value.(KeyPair)
		return true
	})
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

func (m *JWKManager) SignToken(token *jwt.Token) (string, error) {
	jwk := m.GetSignKey()
	if jwk == nil {
		return "", fmt.Errorf("Missing signing key")
	}
	return token.SignedString(jwk.Key)
}

func NewJWKManager() *JWKManager {
	return &JWKManager{keyByKid: make(map[string]KeyPair)}
}
