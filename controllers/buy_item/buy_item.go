package buy_item

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"tbl-backend/database"
	"tbl-backend/models/item"
	// "tbl-backend/models/views"
	// "tbl-backend/services/to_buy_list"
)

var db *sql.DB = database.GetDbConnection()

func GetBuyItems(c *gin.Context) {

	buyItems, err := item.FindItems(db)

	if err != nil {
	 	c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
	 	return
	}

	c.IndentedJSON(http.StatusOK, buyItems)
}

func GetBuyItemById(c *gin.Context) {
	id := c.Param("id")

	res, err := db.Query("SELECT * FROM items")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return;
	}

	var item item.BuyItem
	for res.Next() {
		res.Scan(&item.ID, &item.Name, &item.CurrentQuantity, &item.MinQuantity, &item.SendEmail)

		if (item.ID == id) {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H { "message": "Item not found." })
}

func PostBuyItem(c *gin.Context) {
	var newItem item.BuyItem

	if (c.GetHeader("Content-Type") == "application/json") {
		if err := c.BindJSON(&newItem); err != nil {
			return
		}

		if item := postBuyItem(c, newItem); item != nil {
			c.IndentedJSON(http.StatusCreated, *item)
		}
		return
	} else if (c.GetHeader("Content-Type") == "application/x-www-form-urlencoded") {

		newItem.LoadFromForm(c)

		buyItem := postBuyItem(c, newItem)
		if buyItem == nil {
			return
		}

		c.HTML(http.StatusOK, "form", nil)
		c.HTML(http.StatusOK, "oob-buy-item", buyItem)
		// c.HTML(http.StatusOK, "oob-to-buy-item", buyItem)
		return
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H {
		"message": "Please, use application/json or application/x-www-form-urlencoded as Content-Type value",
	})
}

func PutBuyItem(c *gin.Context) {
	id := c.Param("id")
	var updatedItem item.BuyItem

	if c.GetHeader("Content-Type") == "application/json" {
		if err := c.BindJSON(&updatedItem); err != nil {
			return;
		}

		updatedItem := putBuyItem(c, &updatedItem)

		if updatedItem != nil {
			c.IndentedJSON(http.StatusOK, updatedItem)
			return
		}
	} else if (c.GetHeader("Content-Type") == "application/x-www-form-urlencoded") {
		updatedItem, err := item.FindItem(id, db)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
			return
		}

		updatedItem.LoadFromForm(c)

		updatedItem = putBuyItem(c, updatedItem)

		if updatedItem == nil {
			c.HTML(http.StatusOK, "error-buy-item", nil)
			return
		}

		c.HTML(http.StatusOK, "buy-item", updatedItem)
		return
	}


	c.IndentedJSON(http.StatusNotFound, gin.H { "message": "Item not found." })
}

func DeleteBuyItem(c *gin.Context) {
	id := c.Param("id")

	res, err := db.Query("DELETE FROM items WHERE id=$1 RETURNING *", id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	if c.GetHeader("Content-Type") == "application/json" {
		if res.Next() {
			var deletedItem item.BuyItem
			res.Scan(&deletedItem.ID, &deletedItem.Name, &deletedItem.CurrentQuantity, &deletedItem.MinQuantity, &deletedItem.SendEmail)
			c.IndentedJSON(http.StatusOK, deletedItem)
			return
		}
	} else if c.GetHeader("Content-Type") == "application/x-www-form-urlencoded" {
		c.HTML(http.StatusOK, "", nil)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H { "message": "Item not found." })
	return
}

func postBuyItem(c *gin.Context, newItem item.BuyItem) (*item.BuyItem){
	item, err := newItem.Insert(db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return nil
	}

	return item
}

func putBuyItem(c *gin.Context, updatedItem* item.BuyItem) (*item.BuyItem) {
	updatedItem, err := updatedItem.Update(db)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return nil
	}

	return updatedItem
}
