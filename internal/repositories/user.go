package repositories

import (
	"context"
	"database/sql"
	"fiber-auth-api/internal/types"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type UserRepository struct {
	DB *sql.DB
	log *slog.Logger
}

func NewUserRepository(db *sql.DB, log *slog.Logger) *UserRepository {
	return &UserRepository{
		DB: db,
		log: log,
	}
}

type UserCreateModel struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserCreateDbModel struct {
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

type UserResponseModel struct {
	UserId          string `json:"user_id"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	IsActive        bool   `json:"is_active"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

type Metadata struct {
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	PerPage      int `json:"per_page"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
	TotalCurrentRecords int `json:"total_current_records"`
}

func NewMetadata(totalRecords int) *Metadata {
	if totalRecords == 0 {
		return &Metadata{
			TotalRecords: 0,
			TotalPages:   0,
			CurrentPage:  0,
			PerPage:      0,
			HasNext:      false,
			HasPrev:      false,
			TotalCurrentRecords: 0,
		}
	}

	return &Metadata{
		TotalRecords: totalRecords,
		TotalPages:   totalRecords / 5,
		CurrentPage:  1,
		PerPage:      5,
		HasNext:      totalRecords > 5,
		HasPrev:      false,
		TotalCurrentRecords: 5,
	}
}

type UserListResponseModel struct {
	users []*UserResponseModel
	metadata Metadata
}

func (userRepo UserRepository) CreateUser(user *UserCreateDbModel) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := userRepo.DB.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.IsEmailVerified,
	).Scan(&user.UserId, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if isDuplicateKeyError(err) {
			userRepo.log.Error("User already exists", err)
			return types.ErrDuplicateUser
		}
		userRepo.log.Error("Something went wrong creating user", "error", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (user UserRepository) AuthenticateUser(username string, password string) bool {
	return true
}

func (userRepo UserRepository) GetAllUsers() ([]*UserResponseModel, *Metadata, error) {
	query := `
		SELECT 
			user_id, 
			email, 
			first_name, 
			last_name, 
			username, 
			is_email_verified, 
			is_active,
			count(*) OVER() as total_count 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 
		OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := userRepo.DB.QueryContext(ctx, query, 5, 0)
	if err != nil {
		userRepo.log.Error("Failed to get all users", "error", err)
		return nil, &Metadata{}, err
	}
	defer rows.Close()

	users := make([]*UserResponseModel, 0)

	var totalRecords int
	for rows.Next() {
		user := &UserResponseModel{}
		err := rows.Scan(
			&user.UserId,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.IsActive,
			&user.IsEmailVerified,
			&totalRecords,
		)
		if err != nil {
			userRepo.log.Error("Failed to scan user", "error", err)
			return nil, &Metadata{}, err
		}
		users = append(users, user)
	}

	metadata := NewMetadata(totalRecords)

	return users, metadata, nil
}

func (userRepo UserRepository) FindUserById(userId string) (*UserResponseModel, error) {

	query := `SELECT user_id, email, first_name, last_name, username, is_email_verified, is_active FROM users WHERE user_id = $1`

	var user UserResponseModel

	err := userRepo.DB.QueryRow(query, userId).Scan(
		&user.UserId,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.IsEmailVerified,
		&user.IsActive,
	)

	if err != nil {
		userRepo.log.Error("Failed to get user by id", "error", err)
		return nil, err
	}

	return &user, nil

}

func (user UserRepository) FindUserByUsername(username string) {

}

func (user UserRepository) FindUserByEmail(email string) {}

func (user UserRepository) UpdateUser() {}

func (user UserRepository) UpdateUserById(userId int) {}

func (user UserRepository) UpdateUserPasswordById(userId int) {}

func (user UserRepository) DeleteUser() {}

func (userRepo UserRepository) IsUserExists(email string, username string) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)`

	userExists := false
	err := userRepo.DB.QueryRowContext(context.Background(), query, email, username).Scan(&userExists)
	if err != nil {
		userRepo.log.Error("Failed to check if user exists", "error", err)
		return false, err
	}
	return userExists, nil
}

func isDuplicateKeyError(err error) bool {
    return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}