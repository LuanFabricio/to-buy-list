package token

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(username string) (string, error){
	tokenLifeSpan := 4

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_name"] = username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifeSpan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("TBL_API_SECRET")))
}

func IsTokenValid(tokenString string) bool {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TBL_API_SECRET")), nil
	})

	return err == nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, "")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}
