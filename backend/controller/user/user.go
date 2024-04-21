package user

import (
	"database/sql"
	"net/http"
	"tbl-backend/database"
	"tbl-backend/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB = database.GetDbConnection()

func PostUser(c *gin.Context) {
	var newUser user.UserDTO

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	if !newUser.Hash {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
			return
		}

		newUser.Password = string(hashedPassword)
	}

	user, err := newUser.Insert(db)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	c.IndentedJSON(http.StatusCreated, *user)
}
