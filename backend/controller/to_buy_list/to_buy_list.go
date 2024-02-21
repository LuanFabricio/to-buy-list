package to_buy_list;

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"tbl-backend/item"
	"tbl-backend/controller/buy_item"
)

func GetToBuyList(c *gin.Context) {
	buy_list := make([]item.BuyItem, 0)
	for _, i := range buy_item.BuyItems {
		if i.CurrentQuantity < i.MinQuantity {
			buy_list = append(buy_list, i)
		}
	}

	c.IndentedJSON(http.StatusOK, buy_list)
}
