package api

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	UserId       string
	Username     string
	PasswordHash string
}

type NewUser struct {
	Username string
	Password string
}

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (c *UserRepository) GetUser(username string) (UserModel, error) {
	var user UserModel
	stmt, err := c.db.Prepare("select user_id, username, password_hash from users where username = ?")
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.UserId, &user.Username, &user.PasswordHash)
	return user, err
}

func (c *UserRepository) AuthenticateUser(username string, password string) error {
	user, err := c.GetUser(username)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return errors.New("user password is invalid")
	}
	return nil
}

func (c *UserRepository) AddUser(user NewUser, checkExists bool) error {
	if checkExists {
		_, err := c.GetUser(user.Username)
		if err != sql.ErrNoRows {
			return err
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	_, err := c.db.Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		user.Username, string(hashedPassword),
	)
	return err
}
