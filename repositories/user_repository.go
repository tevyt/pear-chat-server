package repositories

import (
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	RegisterUser(user UserModel) error
}

type UserRepositoryImpl struct {
	dbConnection *sqlx.DB
}

type UserModel struct {
	EmailAddress string `db:"email_address"`
	UserName     string `db:"user_name"`
	PasswordHash string `db:"password_hash"`
	PublicKey    string `db:"public_key"`
}

func NewUserRepository(dbConnection *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{dbConnection: dbConnection}
}

func (userRepository *UserRepositoryImpl) RegisterUser(user UserModel) error {
	registerUserQuery := "INSERT INTO app_user (email_address, user_name, password_hash, public_key) VALUES (:email_address, :user_name, :password_hash, :public_key)"

	_, err := userRepository.dbConnection.NamedExec(registerUserQuery, user)

	return err
}

func (userRepository *UserRepositoryImpl) FindUserByEmailAddress(emailAddress string) (UserModel, error) {
	findUserByEmailAddressQuery := "SELECT * FROM app_user WHERE email_address = :email_address"

	model := UserModel{}

	err := userRepository.dbConnection.Select(&model, findUserByEmailAddressQuery, emailAddress)

	return model, err
}
