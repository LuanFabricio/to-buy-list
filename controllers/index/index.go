package index

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
	buyItems, err := item.FindItems(db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	toBuyItems, err := to_buy_list.FetchToBuyList(db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	vwIndex := views.ViewIndex { BuyItems: buyItems, ToBuyItems: toBuyItems }

	c.HTML(http.StatusOK, "index", vwIndex)
}
