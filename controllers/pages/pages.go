package pages

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"tbl-backend/database"
	"tbl-backend/models/item"
	"tbl-backend/models/views"
	"tbl-backend/services/to_buy_list"
	"tbl-backend/services/token"
)

var db *sql.DB = database.GetDbConnection()

func GetIndex(c *gin.Context) {
	buyItems, err := item.FindItems(db, true)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	toBuyItems, err := to_buy_list.FetchToBuyList(db, 1)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	vwIndex := views.ViewIndex { BuyItems: buyItems, ToBuyItems: toBuyItems }

	c.HTML(http.StatusOK, "index", vwIndex)
}

func GetBuyItemsList(c *gin.Context) {
	buyItems, err := item.FindItems(db, true)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	c.HTML(http.StatusOK, "buy-items-page", buyItems)
}

func GetToBuyItemsList(c *gin.Context) {
	toBuyItems, err := to_buy_list.FetchToBuyList(db, 1)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	c.HTML(http.StatusOK, "to-buy-items-page", toBuyItems)
}

func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", nil)
}

func GetBuyList(c *gin.Context) {
	log.Printf("[INFO]: %v\n", c.Request.Header)
	log.Printf("[INFO]: OK!\n")

	cookie, err := c.Cookie("token")

	if err != nil {
		c.Status(http.StatusUnauthorized)
	}

	log.Printf("[INFO]: %s\n", cookie)
	userId, err := token.ExtractTokenId(cookie)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	log.Printf("[INFO]: User ID: %s\n", userId)
	c.Status(http.StatusOK)
}
