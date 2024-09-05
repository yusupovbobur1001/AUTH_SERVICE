package token

import (
	"auth_service/configs"
	"auth_service/models"
	"auth_service/pkg/logger"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWTToken(request *models.LoginResponse, log logger.ILogger) (*models.Token, error) {
	var (
		tokenValidateTime float64
	)

	accessNewJWt := jwt.New(jwt.SigningMethodHS256)

	accessClaims := accessNewJWt.Claims.(jwt.MapClaims)
	accessClaims["id"] = request.Id
	accessClaims["email"] = request.Email
	accessClaims["user_name"] = request.UserName
	accessClaims["role"] = request.Role
	accessClaims["iat"] = time.Now().Unix()
	accessClaims["exp"] = time.Now().Add(configs.RefreshExpireTime).Unix()

	accessToken, err := accessNewJWt.SignedString(configs.SignKey)
	if err != nil {
		log.Error("this error is signedString-~~~~~~~~~>ERROR", logger.Error(err))
		return nil, err
	}
	refreshNewJWT := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshNewJWT.Claims.(jwt.MapClaims)
	refreshClaims["id"] = request.Id
	refreshClaims["user_name"] = request.UserName
	refreshClaims["role"] = request.Role
	refreshClaims["email"] = request.Email
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["exp"] = time.Now().Add(configs.RefreshExpireTime).Unix()

	tokenValidateTime = float64(refreshClaims["exp"].(int64) - accessClaims["iat"].(int64))
	refreshToken, err := refreshNewJWT.SignedString(configs.SignKey)
	if err != nil {
		log.Error("this error is signed  to string that sign key -~~~~~~~~~ERROR", logger.Error(err))
		return nil, err
	}
	return &models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredTime:  tokenValidateTime,
	}, nil

}

func ExtractClaims(tokenstr string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenstr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(configs.SignKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("faild to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
