package api

import "os"

var config Config

func init() {
	config = Config{}
}

type Config struct {
}

func (c *Config) getJwkPath() string {
	return os.Getenv("JWK_PATH")
}

func (c *Config) getDbPath() string {
	return os.Getenv("DB_PATH")
}

func (c *Config) getExpireTokenSec() uint16 {
	return 1200
}

func (c *Config) getIssuerHost() string {
	return os.Getenv("ISSUER_HOST")
}
