package buy_item

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tbl-backend/database"
	"tbl-backend/models/item"
	"tbl-backend/models/views"
	"tbl-backend/services/to_buy_list"
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
		newItem.ID = "0"
		newItem.Name = c.PostForm("name")

		currentQuantity, _ := strconv.ParseUint(c.PostForm("current_quantity"), 10, 32)
		newItem.CurrentQuantity = uint32(currentQuantity)

		minQuantity, _ := strconv.ParseUint(c.PostForm("min_quantity"), 10, 32)
		newItem.MinQuantity = uint32(minQuantity)

		newItem.SendEmail, _ = strconv.ParseBool(c.PostForm("send_email"))

		if item := postBuyItem(c, newItem); item == nil {
			return
		}

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

		c.HTML(http.StatusOK, "items-list", vwIndex)
		return
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H {
		"message": "Please, use application/json or application/x-www-form-urlencoded as Content-Type value",
	})
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

func postBuyItem(c *gin.Context, newItem item.BuyItem) (*item.BuyItem){
	item, err := newItem.Insert(db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return nil
	}

	return item
}
