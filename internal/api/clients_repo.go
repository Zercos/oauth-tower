package api

import "errors"

type OAuthClient struct {
	ClientId     string
	ClientSecret string
}

type ClientRepository struct {
	db map[string]OAuthClient
}

func (c *ClientRepository) GetClient(clientId string) OAuthClient {
	return c.db[clientId]
}

func (c *ClientRepository) AddClient(client OAuthClient) {
	c.db[client.ClientId] = client
}

func (c *ClientRepository) AuthenticateClient(clientId string, clientSecret string) error {
	client, ok := c.db[clientId]
	if !ok {
		return errors.New("client not exists")
	}
	if client.ClientSecret != clientSecret {
		return errors.New("client secret is invalid")
	}
	return nil
}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{db: make(map[string]OAuthClient)}
}
