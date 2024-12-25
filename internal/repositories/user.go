package repositories

import (
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

type UserCreateResponseModel struct {
	UserId          string    `json:"user_id"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	PasswordHash    string    `json:"password_hash"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	IsActive        bool      `json:"is_active"`
	IsEmailVerified bool      `json:"is_email_verified"`
	LastLoginAt     time.Time `json:"last_login_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}

func (userRepo UserRepository) CreateUser(user *UserCreateResponseModel) error {
	query := `
        INSERT INTO users (
            username, 
            email, 
            password_hash, 
            first_name, 
            last_name, 
            is_active, 
            is_email_verified
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING user_id, created_at, updated_at`

		return userRepo.DB.QueryRowContext(
			context.Background(),
			query,
			user.Username,
			user.Email,
			user.PasswordHash,
			user.FirstName,
			user.LastName,
			user.IsActive,
			user.IsEmailVerified,
		).Scan(&user.UserId, &user.CreatedAt, &user.UpdatedAt)
}

func (user UserRepository) AuthenticateUser(username string, password string) bool {
	return true
}

func (user UserRepository) FindUserById(userId int) {}

func (user UserRepository) FindUserByUsername(username string) {

}

func (user UserRepository) FindUserByEmail(email string) {}

func (user UserRepository) UpdateUser() {}

func (user UserRepository) UpdateUserById(userId int) {}

func (user UserRepository) UpdateUserPasswordById(userId int) {}

func (user UserRepository) DeleteUser() {}

func (user UserRepository) IsUserExists(email string, username string) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)`

	userExists := false
	err := user.DB.QueryRowContext(context.Background(), query, email, username).Scan(&userExists)
	if err != nil {
		// log error
	}
	return userExists, nil
}
