package main

import (
	"tbl-backend/controller/buy_item"
	"tbl-backend/controller/to_buy_list"
	"tbl-backend/controller/user"
	"tbl-backend/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.GetDbConnection()

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/buy_items", buy_item.GetBuyItems)
	router.GET("/buy_items/:id", buy_item.GetBuyItemById)
	router.PUT("/buy_items/:id", buy_item.PutBuyItem)
	router.DELETE("/buy_items/:id", buy_item.DeleteBuyItem)
	router.POST("/buy_items", buy_item.PostBuyItem)
	router.GET("/to_buy_list", to_buy_list.GetToBuyList)
	router.POST("/user", user.PostUser)

	router.Run("localhost:3000")
}
