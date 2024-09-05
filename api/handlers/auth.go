package handlers

import (
	"auth_service/api/email"
	"auth_service/api/token"
	"auth_service/models"
	"auth_service/pkg/logger"
	"auth_service/storage/redis"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Register new user
// @Tags auth
// @Description Register new user
// @Accept  json
// @Produce  json
// @Param  auth body models.CreateUser true "body"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth_service/register [post]
func (h *Handler) Reginster(ctx *gin.Context) {
	request := models.CreateUser{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		h.log.Error("this error is reading from body of user this -~~~~~~~~~~~~~ERROR", logger.Error(err))
		HandleResponse(ctx, h.log, "error is while  reading  by body of   user ", http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.storage.RegisterRepo(ctx, &request)
	if err != nil {
		h.log.Error("this error is register from repository -~~~~~~~~~~~~~ERROR", logger.Error(err))
		HandleResponse(ctx, h.log, "error is while  registering user ", http.StatusInternalServerError, err.Error())
		return
	}
	HandleResponse(ctx, h.log, "SUCCESS", http.StatusOK, response)
	h.log.Error("SUCCESS -~~~~~~~~~~~~~SUCCESS!!!!!!!")
}

// @Summary Login user
// @Tags auth
// @Description Login user
// @Accept  json
// @Produce  json
// @Param body body models.LoginRequest true "body"
// @Success 200 {object} models.Token
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth_service/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	request := models.LoginRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		h.log.Error("this error is reading from body of user this -~~~~~~~~~~~~~ERROR", logger.Error(err))
		HandleResponse(ctx, h.log, "error is while  reading  by body of   user ", http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.storage.Login(ctx, &request)
	if err != nil {
		h.log.Error("this error is register from repository -~~~~~~~~~~~~~ERROR", logger.Error(err))
		HandleResponse(ctx, h.log, "error is while  registering user ", http.StatusInternalServerError, err.Error())
		return
	}

	token, err := token.GenerateJWTToken(response, h.log)
	if err != nil {
		h.log.Error("Token invalite ")
		HandleResponse(ctx, h.log, "token not fount", http.StatusBadRequest, err)
		return
	}
	HandleResponse(ctx, h.log, "SUCCESS", http.StatusOK, token)
}

// @Summary Reset password
// @Tags auth
// @Description Reset password
// @Accept  json
// @Produce  json
// @Param auth body models.ResetPasswordRQ true "body"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth_service/reset_password [post]
func (h *Handler) ResetPassword(ctx *gin.Context) {
	req := models.ResetPasswordRQ{}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.log.Error("Error decoding request body")
		HandleResponse(ctx, h.log, "Error decoding request body", http.StatusBadRequest, err)
		return
	}

	e, err := redis.ReadEmail()
	if err != nil {
		h.log.Error("Error reading email from redis")
		HandleResponse(ctx, h.log, "Error reading email from redis", http.StatusBadRequest, err)
		return
	}

	pass, err := redis.ReadPassword(e)
	if err != nil {
		h.log.Error("Error sending email")
		HandleResponse(ctx, h.log, "Error sending email", http.StatusBadRequest, err)
		return
	}

	if pass != req.ResetToken {
		h.log.Error("Invalid code")
		HandleResponse(ctx, h.log, "Invalid code", http.StatusBadRequest, err)
		return
	}

	resp, err := h.storage.ResetPassword(ctx, &models.ResetPasswordRequest{
		Email:       e,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		h.log.Error("Error updating user password")
		HandleResponse(ctx, h.log, "Error updating user password", http.StatusBadRequest, err)
		return
	}
	fmt.Println(resp)
	HandleResponse(ctx, h.log, "SUCCESS", http.StatusOK, resp)
}

// @Summary Change password
// @Tags auth
// @Description Change password
// @Accept  json
// @Produce  json
// @Param auth body models.ChangePasswordRequest true "body"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth_service/change_password [post]
func (h *Handler) ChangePassword(ctx *gin.Context) {
	req := models.ChangePasswordRequest{}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.log.Error(err.Error())
		HandleResponse(ctx, h.log, "", http.StatusBadRequest, err)
		return
	}

	if len(req.NewPassword) == 0 {
		h.log.Error("New password is empty")
		HandleResponse(ctx, h.log, "New password is empty", http.StatusBadRequest, err)
		return
	}

	if len(req.CurrentPassword) == 0 {
		h.log.Error("Current password is empty")
		HandleResponse(ctx, h.log, "Current password is empty", http.StatusBadRequest, err)
		return
	}

	resp, err := h.storage.ChangePassword(ctx, &req)
	if err != nil {
		h.log.Error(err.Error())
		HandleResponse(ctx, h.log, "not fount", http.StatusBadRequest, err)
		return
	}

	HandleResponse(ctx, h.log, "SUCCESS", http.StatusOK, resp)
}

// @Summary Forgot password
// @Tags auth
// @Description Forgot password
// @Accept  json
// @Produce  json
// @Param auth body models.ForgotPasswordRequest true "body"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /auth_service/forgot_password [post]
func (h *Handler) ForgotPassword(ctx *gin.Context) {
	req := models.ForgotPasswordRequest{}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.log.Error(err.Error())
		HandleResponse(ctx, h.log, err.Error(), http.StatusBadRequest, err)
		return
	}

	code, err := email.Email(req.Email)
	if err != nil {
		h.log.Error(err.Error())
		HandleResponse(ctx, h.log, "not fount", http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(code, "++++++")

	HandleResponse(ctx, h.log, "SUCCESS", http.StatusOK, models.Message{Message: "Password reset instructions sent to your email"})
}
