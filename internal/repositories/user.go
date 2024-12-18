package repositories

import (
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (user UserRepository) CreateUser() {

}

func (user UserRepository) AuthenticateUser(username string, password string) bool {
	return true
}

func (user UserRepository) FindUserById(userId int) {}

func (user UserRepository) FindUserByUsername(username string) {}

func (user UserRepository) FindUserByEmail(email string) {}

func (user UserRepository) UpdateUser() {

}

func (user UserRepository) UpdateUserById(userId int) {}

func (user UserRepository) UpdateUserPasswordById(userId int) {}

func (user UserRepository) DeleteUser() {}
