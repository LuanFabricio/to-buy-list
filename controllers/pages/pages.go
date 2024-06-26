package pages

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	buylist "tbl-backend/models/buy_list"
	"tbl-backend/models/item"
	"tbl-backend/models/user"
	"tbl-backend/models/views"
	"tbl-backend/services/database"
	"tbl-backend/services/logger"
	"tbl-backend/services/to_buy_list"
	"tbl-backend/services/token"
)

var db *sql.DB = database.GetDbConnection()

func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "redirect", gin.H { "PathName": "/buy-list" })
	return

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

func GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register", nil)
}

func GetBuyList(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.HTML(http.StatusOK, "redirect", gin.H { "PathName": "/login" })
		return
	}

	log.Printf("[INFO]: %s\n", cookie)
	userId, err := token.ExtractTokenId(cookie)
	if err != nil {
		c.HTML(http.StatusOK, "redirect", gin.H { "PathName": "/login" })
		return
	}

	user, err := user.FetchUserById(db, userId)
	if err != nil {
		// c.AbortWithError(http.StatusNotFound, err)
		c.HTML(http.StatusOK, "redirect", gin.H { "PathName": "/login" })
		return
	}
	buyListArr := user.FetchBuyLists(db)

	userIdInt, err := strconv.Atoi(userId)
	buyListData := gin.H {
		"UserId": userIdInt,
		"BuyListArr": buyListArr,
	}

	c.HTML(http.StatusOK, "buy_list", buyListData)
}

type BuyListItems struct {
	ListId string;
	Items []item.BuyItem;
	Grid [][]item.BuyItem;
}

func GetBuyListById(c *gin.Context) {
	id := c.Param("id")
	buyListId, _ := strconv.ParseInt(id, 10, 32)

	userToken, err := c.Cookie("token")
	if err != nil {
		logger.Log(logger.ERROR, "%v",err)
		c.HTML(http.StatusUnauthorized, "redirect", gin.H { "PathName": "/login" })
		return
	}

	userID, err := token.ExtractTokenId(userToken)
	if err != nil {
		logger.Log(logger.ERROR, "%v",err)
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

	const MAX_ITEMS_PER_COLUMN = 5
	rows := 1 + (len(buyItemsArr) / MAX_ITEMS_PER_COLUMN)
	buyListItems := BuyListItems {
		ListId: id,
		Items: buyItemsArr,
		Grid: make([][]item.BuyItem, rows),
	}

	for i, buyItem := range buyItemsArr {
		colIdx := i % MAX_ITEMS_PER_COLUMN
		rowIdx := i / MAX_ITEMS_PER_COLUMN

		if len(buyListItems.Grid[rowIdx]) == 0 {
			buyListItems.Grid[rowIdx] = make([]item.BuyItem, MAX_ITEMS_PER_COLUMN)
		}
		buyListItems.Grid[rowIdx][colIdx] = buyItem
	}

	logger.Log(logger.INFO, "grid: %v", buyListItems.Grid)
	c.HTML(http.StatusOK, "buy-items-page", buyListItems)
}

func GetModal(c *gin.Context) {
	modalType := "modal-" + c.Param("id")

	modalItems := gin.H {}
	queries := c.Request.URL.Query()
	for k := range queries {
		modalItems[k] = queries.Get(k)
	}

	log.Println(modalType)

	c.HTML(http.StatusOK, modalType, modalItems)
}

func GetRedirect(c *gin.Context) {
	path_name := c.Param("path_name")

	path := "/" + path_name
	c.Header("HX-Redirect", path)
	c.Status(http.StatusOK)
}
