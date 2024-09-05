package handlers

import (
	"auth_service/models"
	"auth_service/pkg/logger"
	"auth_service/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage storage.IAuthStorage
	log     logger.ILogger
}

func NewHandler(storage storage.IAuthStorage, log logger.ILogger) *Handler {
	return &Handler{
		storage: storage,
		log:     log,
	}
}

func HandleResponse(c *gin.Context, log logger.ILogger, msg string, statusCode int, data interface{}) {

	resp := models.Response{}

	switch code := statusCode; {
	case code < 400:
		resp.Description = "OK"
		log.Info("~~~~> OK", logger.String("msg", msg), logger.Any("status", code))
	case code == 401:
		resp.Description = "Unauth_serviceorized"
		log.Error("???? Unauth_serviceorized", logger.String("msg", msg), logger.Any("status", code))
	case code < 500:
		resp.Description = "Bad Request"
		log.Error("!!!!! BAD REQUEST", logger.String("msg", msg), logger.Any("status", code))
	default:
		resp.Description = "Internal Server Error"
		log.Error("!!!!! INTERNAL SERVER ERROR", logger.String("msg", msg), logger.Any("status", code), logger.Any("error", data))
	}

	resp.StatusCode = statusCode
	resp.Data = data

	c.JSON(resp.StatusCode, resp)
}
