package user

import (
	"database/sql"
	"log"
	"net/http"
	"tbl-backend/database"
	"tbl-backend/models/user"

	"github.com/gin-gonic/gin"
	"crypto/sha256"
	"encoding/hex"
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

	userHash, err := fetchUserHash(login.Username)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		c.Status(http.StatusNotFound)
		return
	}
	log.Printf("Current: %s : %s (%s)", login.Username, login.Password, login.Username + login.Password)

	log.Printf("User hash: %s", userHash)
	currentHash := string(currentHashBytes)
	log.Printf("Current hash: %s", currentHash)
	if userHash == currentHash {
		c.Header("HX-Redirect", "/")
		// TODO(Luan): Create auth token
		c.Header("Set-Cookie", "token=cookie123")
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusUnauthorized)
}

func fetchUserHash(username string) (string, error) {
	var user user.User

	if err := user.FindByEmail(db, username); err != nil {
		return "", err
	}

	return user.Password, nil
}

func generateUserHash(email, password string) string {
	originalString := email + password

	sha := sha256.New()
	sha.Write([]byte(originalString))
	return hex.EncodeToString(sha.Sum(nil))
}
