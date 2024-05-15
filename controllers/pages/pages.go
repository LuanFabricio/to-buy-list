package pages

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"tbl-backend/database"
	"tbl-backend/models/item"
	"tbl-backend/models/views"
	"tbl-backend/services/to_buy_list"
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
