package buy_item

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"tbl-backend/item"
	"tbl-backend/database"
)

var db *sql.DB = database.GetDbConnection()

func GetBuyItems(c *gin.Context) {
	res, err := db.Query("SELECT * FROM items")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	buyItems := []item.BuyItem{}
	for res.Next() {
		var b item.BuyItem
		res.Scan(&b.ID, &b.Name, &b.CurrentQuantity, &b.MinQuantity, &b.SendEmail)
		buyItems = append(buyItems, b)
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

	if err := c.BindJSON(&newItem); err != nil {
		return
	}

	item, err := newItem.Insert(db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	c.IndentedJSON(http.StatusCreated, *item)
}

func PutBuyItem(c *gin.Context) {
	id := c.Param("id")
	var updatedItem item.BuyItem

	if err := c.BindJSON(&updatedItem); err != nil {
		return;
	}

	res, err := db.Query(`
		UPDATE items
		SET name = $1, current_quantity = $2, min_quantity = $3, send_email = $4
		WHERE id = $5
		RETURNING *`,
		updatedItem.Name, updatedItem.CurrentQuantity,
		updatedItem.MinQuantity, updatedItem.SendEmail,
		id);

	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return;
	}

	if res.Next() {
		var newItem item.BuyItem
		res.Scan(&newItem.ID, &newItem.Name, &newItem.CurrentQuantity, &newItem.MinQuantity, &newItem.SendEmail)
		c.IndentedJSON(http.StatusOK, updatedItem)
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

	if res.Next() {
		var deletedItem item.BuyItem
		res.Scan(&deletedItem.ID, &deletedItem.Name, &deletedItem.CurrentQuantity, &deletedItem.MinQuantity, &deletedItem.SendEmail)
		c.IndentedJSON(http.StatusOK, deletedItem)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H { "message": "Item not found." })
}
