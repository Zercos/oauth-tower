package api

import (
	"errors"
)

type OAuthClient struct {
	ClientId     string
	ClientSecret string
	RedirectURI  string
}

type ClientRepository struct {
	db *DB
}

func (c *ClientRepository) GetClient(clientId string) (OAuthClient, error) {
	var client OAuthClient
	stmt, err := c.db.Prepare("select client_id, client_secret, redirect_uri from clients where client_id = ?")
	if err != nil {
		return client, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(clientId).Scan(&client.ClientId, &client.ClientSecret, &client.RedirectURI)
	if err != nil {
		return client, err
	}
	return client, nil
}

func (c *ClientRepository) AddClient(client OAuthClient) error {
	_, err := c.db.Exec(
		"INSERT INTO clients (client_id, client_secret, redirect_uri)  VALUES (?, ?, ?)",
		client.ClientId, client.ClientSecret, client.RedirectURI,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientRepository) AuthenticateClient(clientId string, clientSecret string) error {
	client, err := c.GetClient(clientId)
	if err != nil {
		return err
	}
	if client.ClientSecret != clientSecret {
		return errors.New("client secret is invalid")
	}
	return nil
}

func NewClientRepository(db *DB) *ClientRepository {
	return &ClientRepository{db: db}
}
