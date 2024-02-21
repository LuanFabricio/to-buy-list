package buy_item

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"tbl-backend/item"
)

var BuyItems = []item.BuyItem {
	{ ID: "32dsa", Name: "T1", CurrentQuantity: 1, MinQuantity: 1, SendEmail: true },
	{ ID: "3-As2", Name: "T2", CurrentQuantity: 2, MinQuantity: 1, SendEmail: false },
	{ ID: "sa32d", Name: "T3", CurrentQuantity: 3, MinQuantity: 1, SendEmail: true },
}

func GetBuyItems(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, BuyItems)
}

func GetBuyItemById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range BuyItems {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return;
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H { "message": "Item not found." })
}

func PostBuyItem(c *gin.Context) {
	var newItem item.BuyItem

	if err := c.BindJSON(&newItem); err != nil {
		return
	}

	newItem.ID = "123"
	BuyItems = append(BuyItems, newItem)
	c.IndentedJSON(http.StatusCreated, newItem)
}

func PutBuyItem(c *gin.Context) {
	id := c.Param("id")
	var updatedItem item.BuyItem

	if err := c.BindJSON(&updatedItem); err != nil {
		return;
	}

	for idx, a := range BuyItems {
		if a.ID == id {
			BuyItems[idx] = updatedItem
			c.IndentedJSON(http.StatusOK, BuyItems[idx])
			return;
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H { "message": "Item not found." })
}

func DeleteBuyItem(c *gin.Context) {
	id := c.Param("id")

	for idx, a := range BuyItems {
		if a.ID == id {
			BuyItems = append(BuyItems[:idx], BuyItems[idx+1:]...)
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H { "message": "Item not found." })
}
