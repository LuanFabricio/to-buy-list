package to_buy_list

import (
	"database/sql"
	"net/http"
	"tbl-backend/database"
	"tbl-backend/item"

	"github.com/gin-gonic/gin"
)

var db *sql.DB = database.GetDbConnection()

func GetToBuyList(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM items")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	var buy_list []item.BuyItem = []item.BuyItem{}
	var buy_item item.BuyItem
	for rows.Next(){
		rows.Scan(&buy_item.ID, &buy_item.Name, &buy_item.CurrentQuantity, &buy_item.MinQuantity, &buy_item.SendEmail)
		if buy_item.CurrentQuantity < buy_item.MinQuantity {
			buy_list = append(buy_list, buy_item)
		}
	}

	c.IndentedJSON(http.StatusOK, buy_list)
}
