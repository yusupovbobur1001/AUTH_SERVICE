package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"auth_service/configs"
	"auth_service/models"
	"auth_service/pkg/logger"

	"github.com/stretchr/testify/assert"
)


type MockDB struct {
	log logger.ILogger
	ctx context.Context
	cfg configs.Config
}

func TestAuthRepository_RegisterRepo(t *testing.T) {
	mockDB := &MockDB{} 
	db, err := NewStore(mockDB.ctx, mockDB.log, mockDB.cfg)
	if err != nil {
		fmt.Println(err, "+++++++++++++++++++++")
		panic(err)
	}
	defer db.Close()

	repo := NewAuthRepository(db.DB, mockDB.log)

	request := &models.CreateUser{
		UserName:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "123456",
		Role:         "user",
	}

	expectedUser := &models.User{
		Id:           "3758dbeb-a8a0-4847-9d34-7f5db7f1523c",
		UserName:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	user, err := repo.RegisterRepo(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func TestAuthRepository_Login(t *testing.T) {
	mockDB := &MockDB{}
	
	db, err := NewStore(mockDB.ctx, mockDB.log, mockDB.cfg)
	if err != nil {
		fmt.Println(err, "+++++++++++++++++++++")
		panic(err)
	}
	defer db.Close()

	repo := NewAuthRepository(db.DB, mockDB.log)

	request := &models.LoginRequest{
		Email:        "test@example.com",
		PasswordHash: "123456",
	}

	expectedResponse := &models.LoginResponse{
		Id:       "67f5e17d-04cf-438f-b322-d8da6fed3281",
		UserName: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}

	response, err := repo.Login(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Email, response.Email)
}

func TestAuthRepository_ChangePassword(t *testing.T) {
	mockDB := &MockDB{}

	db, err := NewStore(mockDB.ctx, mockDB.log, mockDB.cfg)
	if err != nil {
		fmt.Println(err, "+++++++++++++++++++++")
		panic(err)
	}

	repo := NewAuthRepository(db.DB, mockDB.log)

	request := &models.ChangePasswordRequest{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword",
	}

	expectedMessage := &models.Message{Message: "Password changed successfully"}

	message, err := repo.ChangePassword(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessage.Message, message.Message)
}

func TestAuthRepository_ResetPassword(t *testing.T) {
	mockDB := &MockDB{}

	db, err := NewStore(mockDB.ctx, mockDB.log, mockDB.cfg)
	if err != nil {
		fmt.Println(err, "+++++++++++++++++++++")
		panic(err)
	}

	repo := NewAuthRepository(db.DB, mockDB.log)

	request := &models.ResetPasswordRequest{
		Email:       "test@example.com",
		NewPassword: "newpassword",
	}

	expectedMessage := &models.Message{Message: "Password reset successfully"}

	message, err := repo.ResetPassword(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessage.Message, message.Message)
}
