package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"

	"tbl-backend/controllers/buy_item"
	"tbl-backend/controllers/pages"
	"tbl-backend/controllers/user"
	"tbl-backend/database"
	"tbl-backend/services/to_buy_list"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = to_buy_list.SendToBuyListEmail(database.GetDbConnection())

	if err != nil {
		log.Fatal(err)
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
		return
	}

	job, err := scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(12, 00, 00),
				gocron.NewAtTime(13, 12, 00),
			),
		),
		gocron.NewTask(
			to_buy_list.SendToBuyListEmail,
			database.GetDbConnection(),
		),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Job ID: %v", job.ID())
	scheduler.Start()
	defer scheduler.Shutdown()

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
	// router.GET("/to_buy_list", to_buy_list.GetToBuyList)
	router.POST("/user", user.PostUser)

	router.Run("localhost:3000")
}
