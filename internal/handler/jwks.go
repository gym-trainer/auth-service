package handler

import (
	"crypto/rsa"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gym-trainer/auth-service/internal/token"
)

type JWKSHandler struct {
	jwks map[string]interface{}
}

func NewJWKSHandler(pubKey *rsa.PublicKey) *JWKSHandler {
	return &JWKSHandler{
		jwks: token.GenerateJWKS(pubKey),
	}
}

func (h *JWKSHandler) ServeJWKS(c *gin.Context) {
	c.JSON(http.StatusOK, h.jwks)
}
