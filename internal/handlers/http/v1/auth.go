package apiv1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 24 * time.Hour

type AuthHandler struct {
	adminLogin    string
	adminPassword string
	jwtSecret     string
}

func NewAuthHandler(adminLogin, adminPassword, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		adminLogin:    adminLogin,
		adminPassword: adminPassword,
		jwtSecret:     jwtSecret,
	}
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})

		return
	}

	if req.Login != h.adminLogin || req.Password != h.adminPassword {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": req.Login,
		"exp":   time.Now().Add(tokenTTL).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})

		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}
