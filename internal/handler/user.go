package handler

import (
	"errors"
	"net/http"
	"strings"

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

	result, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, model.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": model.ErrInternalServer.Error(),
		})
		return
	}

	setRefreshCookie(c, result.RefreshToken)
	c.JSON(http.StatusCreated, model.AuthResponse{
		AccessToken: result.AccessToken,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input model.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": model.ErrInternalServer.Error(),
		})
		return
	}

	setRefreshCookie(c, result.RefreshToken)
	c.JSON(http.StatusOK, model.AuthResponse{
		AccessToken: result.AccessToken,
	})
}

func (h *UserHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": model.ErrUnauthorized.Error(),
		})
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, model.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": model.ErrInternalServer.Error(),
		})
		return
	}

	setRefreshCookie(c, result.RefreshToken)
	c.JSON(http.StatusOK, model.AuthResponse{
		AccessToken: result.AccessToken,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")

	var accessToken string
	authHeader := c.GetHeader("Authorization")

	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		accessToken = authHeader[7:]
	}

	err := h.service.Logout(c.Request.Context(), accessToken, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": model.ErrInternalServer.Error()})
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
