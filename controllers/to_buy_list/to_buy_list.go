package to_buy_list

import (
	"database/sql"
	"net/http"
	"tbl-backend/database"
	"tbl-backend/services/to_buy_list"

	"github.com/gin-gonic/gin"
)

var db *sql.DB = database.GetDbConnection()

func GetToBuyList(c *gin.Context) {
	buy_list, err := to_buy_list.FetchToBuyList(db)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	c.IndentedJSON(http.StatusOK, buy_list)
}
