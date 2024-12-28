package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"trackonomy/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header is missing", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Authorization header format is invalid", nil)
			c.Abort()
			return
		}

		tokenStr := parts[1]
		userID, err := validateToken(tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}

		// Set userID in context for later retrieval
		c.Set("userID", userID)
		c.Next()
	}
}

func validateToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user_id from the claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("token claims are invalid")
		}
		return uint(userIDFloat), nil
	}

	return 0, errors.New("token is invalid")
}
