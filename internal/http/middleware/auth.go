package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/underground20/sso-jwt-token/pkg/jwt/token"
	"github.com/underground20/sso-jwt-token/pkg/jwt/user"
)

var AuthKey struct{}

type AuthContext struct {
	Scopes []string
	Roles  []string
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		var claims user.Claims
		_, err := token.Parse(tokenString, secret, &claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		authContext := AuthContext{
			Scopes: claims.Scopes,
			Roles:  claims.Roles,
		}

		ctx := context.WithValue(c.Request.Context(), AuthKey, &authContext)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
