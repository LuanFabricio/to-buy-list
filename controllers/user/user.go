package user

import (
	"database/sql"
	"log"
	"net/http"
	"tbl-backend/database"
	"tbl-backend/models/user"

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


	log.Printf("Login: %v\n", login)
	// TODO(Luan): Check user credentials
	if len(login.Username) == len(login.Password) {
		c.Header("HX-Redirect", "/")
		// TODO(Luan): Create auth token
		c.Header("Set-Cookie", "token=cookie123")
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusUnauthorized)
}
