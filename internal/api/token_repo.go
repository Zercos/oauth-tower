package api

type TokenRepository struct {
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func (t *TokenRepository) GetUserIdForToken(authToken string) (string, error) {
	return "", nil

}

func (t *TokenRepository) RemoveToken(authToken string) error {
	return nil
}

func (t *TokenRepository) SetAuthToken(userId string, authToken string) error {
	return nil
}
