package main

import (
	"tbl-backend/controllers/buy_item"
	"tbl-backend/controllers/index"
	"tbl-backend/controllers/to_buy_list"
	"tbl-backend/controllers/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())
	router.LoadHTMLGlob("views/*.html")

	router.GET("/", index.GetIndex)

	router.GET("/buy_items", buy_item.GetBuyItems)
	router.GET("/buy_items/:id", buy_item.GetBuyItemById)
	router.PUT("/buy_items/:id", buy_item.PutBuyItem)
	router.DELETE("/buy_items/:id", buy_item.DeleteBuyItem)
	router.POST("/buy_items", buy_item.PostBuyItem)
	router.GET("/to_buy_list", to_buy_list.GetToBuyList)
	router.POST("/user", user.PostUser)

	router.Run("localhost:3000")
}
