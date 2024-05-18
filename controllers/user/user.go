package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"tbl-backend/database"
	"tbl-backend/models/user"
	"tbl-backend/services/token"

	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

var db *sql.DB = database.GetDbConnection()

func PostUser(c *gin.Context) {
	var newUser user.UserDTO

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	if !newUser.Hash {
		hashedPassword := generateUserHash(newUser.Username, newUser.Password)

		newUser.Password = string(hashedPassword)
	}

	user, err := newUser.Insert(db)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	c.IndentedJSON(http.StatusCreated, *user)
}

// TODO(Luan): Move struct to models or refactor UserDTO
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func AuthUser(c *gin.Context) {
	var login Login

	if err := c.ShouldBind(&login); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	currentHashBytes := generateUserHash(login.Username, login.Password)

	user, err := fetchUser(login.Username)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		c.Status(http.StatusNotFound)
		return
	}
	log.Printf("Current: %s : %s (%s)", login.Username, login.Password, login.Username + login.Password)

	log.Printf("User hash: %s", user.Password)
	currentHash := string(currentHashBytes)
	log.Printf("Current hash: %s", currentHash)
	if user.Password == currentHash {
		c.Header("HX-Redirect", "/")
		log.Printf("[INFO]: userID=%s", user.ID)
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
