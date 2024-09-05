package postgres

import (
	"auth_service/models"
	"auth_service/pkg/helper"
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	DB  *pgxpool.Pool
	log logger.ILogger
}

func NewAuthRepository(db *pgxpool.Pool, log logger.ILogger) storage.IAuthStorage {
	return &AuthRepository{
		DB:  db,
		log: log,
	}
}

func (repo *AuthRepository) RegisterRepo(ctx context.Context, request *models.CreateUser) (*models.User, error) {
	var (
		err        error
		response   models.User
		createTime = time.Now()
		password   string
	)
	password = helper.HashPassword(request.PasswordHash)
	query := `INSERT INTO 
    users (
           user_name,
           email,
           password_hash,
           role,
           created_at)
    VALUES($1,$2,$3,$4,$5) 
    returning 
    id,
    user_name,
    email,
    password_hash,
    role,
    created_at,
    updated_at `

	err = repo.DB.QueryRow(ctx, query,
		request.UserName,
		request.Email,
		password,
		request.Role,
		createTime,
	).Scan(
		&response.Id,
		&response.UserName,
		&response.Email,
		&response.PasswordHash,
		&response.Role,
		&response.CreatedAt,
		&response.UpdatedAt,
	)
	if err != nil {
		repo.log.Error("this error is register user it can used that ERROR-~~~~>", logger.Error(err))
		return nil, err
	}
	return &response, nil
}

func (repo *AuthRepository) Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error) {
	var (
		err      error
		query    string
		response models.LoginResponse
		password string
	)
	password = helper.HashPassword(request.PasswordHash)

	query = `
	SELECT
    	id,
    	user_name,
    	email,
    	role
    FROM 
		users  
	WHERE 
		email=$1 
	AND 
		password_hash=$2 `

	err = repo.DB.QueryRow(ctx, query, request.Email, password).
		Scan(&response.Id,
			&response.UserName,
			&response.Email,
			&response.Role,
		)
	if err != nil {
		repo.log.Error("this error is ")
		return nil, err
	}
	return &response, nil
}

func (repo *AuthRepository) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) (*models.Message, error) {
	passwordHash := helper.HashPassword(req.CurrentPassword)
	newPassword := helper.HashPassword(req.NewPassword)
	q := `
	update 
		users 
	set 
		password_hash=$1, 
		updated_at=$2 
	where 
		password_hash=$3
	`
	_, err := repo.DB.Exec(ctx, q, newPassword, time.Now(), passwordHash)
	if err != nil {
		return nil, err
	}
	return &models.Message{Message: "Password changed successfully"}, nil
}

func (repo *AuthRepository) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.Message, error) {
	fmt.Println(req.Email, req.ResetToken)
	newPassword := helper.HashPassword(req.NewPassword)
	q := `
	update 
		users
	set 
		password_hash=$1,
		updated_at=$2
	where 
		email=$3
	`
	_, err := repo.DB.Exec(ctx, q, newPassword, time.Now(), req.Email)
	if err != nil {
		return nil, err
	}

	return &models.Message{Message: "Password reset successfully"}, nil
}
