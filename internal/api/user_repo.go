package api

import "errors"

type UserModel struct {
	UserId       string
	Username     string
	PasswordHash string
}

type UserRepository struct {
	db *DB
}

func (c *UserRepository) GetUser(username string) (UserModel, error) {
	var user UserModel
	stmt, err := c.db.Prepare("select user_id, username, password_hash from users where username = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.UserId, &user.Username, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (c *UserRepository) AuthenticateUser(username string, password string) error {
	user, err := c.GetUser(username)
	if err != nil {
		return err
	}
	if user.PasswordHash != password {
		return errors.New("user password is invalid")
	}
	return nil
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}
