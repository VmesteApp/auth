package middlewares

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	jwt.RegisteredClaims
	Uid  uint64
	Role string
}

var jwtPattern = regexp.MustCompile(`^Bearer\s([A-Za-z0-9\-._~+\/]+=*)$`)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Access denied. No token provided."})
			c.Abort()
			return
		}

		matches := jwtPattern.FindStringSubmatch(header)

		if len(matches) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format."})
			c.Abort()
			return
		}

		var userClaim UserClaim

		token, err := jwt.ParseWithClaims(matches[1], &userClaim, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Expired token."})
			c.Abort()
			return
		}
		if err != nil || !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token."})
			c.Abort()
			return
		}

		c.Set("uid", userClaim.Uid)
		c.Set("role", userClaim.Role)

		c.Next()
	}
}
