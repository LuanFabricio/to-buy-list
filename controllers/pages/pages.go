package pages

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tbl-backend/database"
	buylist "tbl-backend/models/buy_list"
	"tbl-backend/models/item"
	"tbl-backend/models/user"
	"tbl-backend/models/views"
	"tbl-backend/services/to_buy_list"
	"tbl-backend/services/token"
)

var db *sql.DB = database.GetDbConnection()

func GetIndex(c *gin.Context) {
	buyItems, err := item.FindItems(db, true)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	toBuyItems, err := to_buy_list.FetchToBuyList(db, 1)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	vwIndex := views.ViewIndex { BuyItems: buyItems, ToBuyItems: toBuyItems }

	c.HTML(http.StatusOK, "index", vwIndex)
}

func GetBuyItemsList(c *gin.Context) {
	buyItems, err := item.FindItems(db, true)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	c.HTML(http.StatusOK, "buy-items-page", buyItems)
}

func GetToBuyItemsList(c *gin.Context) {
	toBuyItems, err := to_buy_list.FetchToBuyList(db, 1)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H { "message": err })
		return
	}

	c.HTML(http.StatusOK, "to-buy-items-page", toBuyItems)
}

func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", nil)
}

func GetBuyList(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.HTML(http.StatusUnauthorized, "redirect", gin.H { "PathName": "/login" })
		return
	}

	log.Printf("[INFO]: %s\n", cookie)
	userId, err := token.ExtractTokenId(cookie)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "redirect", gin.H { "PathName": "/login" })
		return
	}

	user, err := user.FetchUserById(db, userId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	buyListArr := user.FetchBuyLists(db)

	c.HTML(http.StatusOK, "buy_list", buyListArr)
}

func GetBuyListById(c *gin.Context) {
	id := c.Param("id")
	buyListId, _ := strconv.ParseInt(id, 10, 32)

	userToken, err := c.Cookie("token")
	if err != nil {
		log.Printf("[ERROR] %v", err)
		c.HTML(http.StatusUnauthorized, "redirect", gin.H { "PathName": "/login" })
		return
	}

	userID, err := token.ExtractTokenId(userToken)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		c.HTML(http.StatusUnauthorized, "redirect", gin.H { "PathName": "/login" })
		return
	}
	buyList := buylist.BuyList{ ID: int(buyListId) }

	haveAccess, err := buyList.UserHaveAccess(db, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !haveAccess {
		c.HTML(http.StatusUnauthorized, "redirect", gin.H { "PathName": "/buy-list" })
		return
	}

	buyItemsArr, err := buyList.FetchItems(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(buyList.ID)
	log.Println(buyItemsArr)

	c.HTML(http.StatusOK, "buy-items-page", buyItemsArr)
}
