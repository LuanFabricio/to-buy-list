package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"tbl-backend/database"
	"tbl-backend/models/user"
	"tbl-backend/services/logger"
	"tbl-backend/services/token"

	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

var db *sql.DB = database.GetDbConnection()

func PostUser(c *gin.Context) {
	var newUser user.UserDTO

	if err := c.ShouldBind(&newUser); err != nil {
		return
	}

	if newUser.Password != newUser.ConfirmPassword {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": "Wrong password confirmation" })
		return
	}

	currentHashBytes := generateUserHash(newUser.Username, newUser.Password)
	newUser.Password = currentHashBytes
	user, err := newUser.Insert(db)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	if c.ContentType() == "application/x-www-form-urlencoded" {
		c.Header("HX-Redirect", "/login")
		c.Status(http.StatusOK)
		return
	}

	c.IndentedJSON(http.StatusCreated, *user)
}

func AuthUser(c *gin.Context) {
	var login user.User

	if err := c.ShouldBind(&login); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	currentHashBytes := generateUserHash(login.Username, login.Password)

	user, err := fetchUser(login.Username)
	if err != nil {
		logger.Log(logger.ERROR, "%v", err)
		c.Status(http.StatusNotFound)
		return
	}
	logger.Log(logger.INFO, "Current: %s : %s (%s)", login.Username, login.Password, login.Username + login.Password)

	logger.Log(logger.INFO, "User hash: %s", user.Password)
	currentHash := string(currentHashBytes)
	logger.Log(logger.INFO, "Current hash: %s", currentHash)
	if user.Password == currentHash {
		c.Header("HX-Redirect", "/")
		logger.Log(logger.INFO, "userID=%s", user.ID)
		authToken, err := token.GenerateToken(user.ID)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.Header("Set-Cookie", fmt.Sprintf("token=%s", authToken))
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusUnauthorized)
}

func fetchUser(username string) (user.User, error) {
	var user user.User

	if err := user.FindByEmail(db, username); err != nil {
		return user, err
	}

	return user, nil
}

func generateUserHash(email, password string) string {
	originalString := email + password

	sha := sha256.New()
	sha.Write([]byte(originalString))
	return hex.EncodeToString(sha.Sum(nil))
}
