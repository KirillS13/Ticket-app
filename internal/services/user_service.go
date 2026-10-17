package services

import (
	"Ticket-app/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (service *Service) CreateUser(email string, username string, password string, ctx context.Context) error {
	var user domain.User
	passkey, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user = domain.User{
		Email:     email,
		Username:  username,
		Password:  passkey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := "INSERT INTO users (email, username, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err = service.db.Exec(ctx, req, user.Email, username, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) Login(email string, password string, secretKey string, ctx context.Context) (string, error) {
	var user domain.User
	req := "SELECT id, username, password_hash FROM users WHERE email = $1"
	query := service.db.QueryRow(ctx, req, email)
	err := query.Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return "", fmt.Errorf("error canning user: %w", err)
	}
	isCorrect := VerifyPassword(user.Password, password)
	if !isCorrect {
		return "", fmt.Errorf("invalid password")
	}
	jwt, err := GenerateToken(email, user.Id, secretKey)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}
	req = "UPDATE users SET updated_at = $1 WHERE email = $2"
	_, err = service.db.Exec(ctx, req, time.Now(), email)
	if err != nil {
		return "", fmt.Errorf("error updating user: %w", err)
	}

	return jwt, nil
}
