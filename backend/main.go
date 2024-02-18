package main

import (
	// "net/http"
	//"tbl-backend/item"
	"tbl-backend/controller/buy_item"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/buy_items", buy_item.GetBuyItems)
	router.GET("/buy_items/:id", buy_item.GetBuyItemById)
	router.PUT("/buy_items/:id", buy_item.PutBuyItem)
	router.DELETE("/buy_items/:id", buy_item.DeleteBuyItem)
	router.POST("/buy_items", buy_item.PostBuyItem)

	router.Run("localhost:3000")
}
