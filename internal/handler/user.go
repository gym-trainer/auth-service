package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gym-trainer/auth-service/internal/model"
	"github.com/gym-trainer/auth-service/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func setRefreshCookie(c *gin.Context, token string) {
	maxAge := 7 * 24 * 60 * 60

	c.SetCookie("refresh_token", token, maxAge, "/", "", false, true)
}

func (h *UserHandler) Register(c *gin.Context) {
	var input model.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Register(c, input)
	if err != nil {
		if errors.Is(err, model.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": model.ErrInternalServer.Error()})
		return
	}

	setRefreshCookie(c, result.RefreshToken)
	c.JSON(http.StatusCreated, model.AuthResponse{
		User:        result.User,
		AccessToken: result.AccessToken,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input model.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Login(c, input)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": model.ErrInternalServer.Error()})
	}

	setRefreshCookie(c, result.RefreshToken)
	c.JSON(http.StatusOK, model.AuthResponse{
		User:        result.User,
		AccessToken: result.AccessToken,
	})
}
