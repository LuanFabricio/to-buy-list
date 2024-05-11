package main

import (
	"tbl-backend/controllers/buy_item"
	"tbl-backend/controllers/pages"
	"tbl-backend/controllers/to_buy_list"
	"tbl-backend/controllers/user"
	"log"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	router := gin.Default()

	router.Static("/static", "./views/static/")

	router.Use(cors.Default())
	router.LoadHTMLGlob("views/*.html")

	// HTML pages
	router.GET("/", pages.GetIndex)
	router.GET("/buy-items", pages.GetBuyItemsList)
	router.GET("/to-buy-items", pages.GetToBuyItemsList)

	// Endpoints
	router.GET("/buy_items", buy_item.GetBuyItems)
	router.GET("/buy_items/:id", buy_item.GetBuyItemById)
	router.PUT("/buy_items/:id", buy_item.PutBuyItem)
	router.DELETE("/buy_items/:id", buy_item.DeleteBuyItem)
	router.POST("/buy_items", buy_item.PostBuyItem)
	router.GET("/to_buy_list", to_buy_list.GetToBuyList)
	router.POST("/user", user.PostUser)

	router.Run("localhost:3000")
}
