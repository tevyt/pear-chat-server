package repositories

import (
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	dbConnection *sqlx.DB
}

type UserModel struct {
	EmailAddress string `db:"email_address"`
	UserName     string `db:"user_name"`
	PasswordHash string `db:"password_hash"`
	PublicKey    string `db:"public_key"`
}

func NewUserRepository(dbConnection *sqlx.DB) *UserRepository {
	return &UserRepository{dbConnection: dbConnection}
}

func (userRepository *UserRepository) RegisterUser(user *UserModel) error {
	registerUserQuery := "INSERT INTO app_user (email_address, user_name, password_hash, public_key) VALUES (:email_address, :user_name, :password_hash, :public_key)"

	_, err := userRepository.dbConnection.NamedExec(registerUserQuery, user)

	return err
}
