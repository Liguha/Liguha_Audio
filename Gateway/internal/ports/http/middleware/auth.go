package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userIDHeader        = "x-user-id"
	userRoleHeader      = "x-user-role"
)

func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(authorizationHeader)
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(jwtSecret), nil
		})
		if err != nil {
			log.Error(err)
			c.Status(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Request.Header.Del(authorizationHeader)
		c.Request.Header.Add(userIDHeader, userID)
		c.Request.Header.Add(userRoleHeader, role)
		c.Next()
	}
}
