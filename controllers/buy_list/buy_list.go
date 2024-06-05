package buy_list

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tbl-backend/database"
	buylist "tbl-backend/models/buy_list"
	"tbl-backend/services/logger"
	"tbl-backend/services/token"
)

var db *sql.DB = database.GetDbConnection()

func PostBuyList(c *gin.Context) {
	var buyList buylist.BuyList

	userToken, err := c.Cookie("token")
	if err != nil  {
		c.HTML(http.StatusOK, "modal-error", gin.H { "error": err })
		return
	}

	userIdStr, err := token.ExtractTokenId(userToken)
	logger.Log(logger.INFO, "User ID: %s", userIdStr)
	if err != nil  {
		c.HTML(http.StatusOK, "modal-error", gin.H { "error": err })
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		c.HTML(http.StatusOK, "modal-error", gin.H { "error": err })
		return
	}

	if err := c.ShouldBind(&buyList); err != nil {
		c.HTML(http.StatusOK, "modal-error", gin.H { "error": err })
		return
	}

	buyList.OwnerUserID = int(userId)
	logger.Log(logger.INFO, "BuyList: %v", buyList)
	_, err = buyList.Insert(db)
	if err != nil {
		c.HTML(http.StatusOK, "modal-error", gin.H { "error": err })
		return
	}

//	succesMap := gin.H {
//		"Title": "Success!",
//		"success": fmt.Sprintf("Buy list %s created with success!", newBuyList.Name),
//	}
//	c.HTML(http.StatusOK, "modal-success", successMap)

	// TODO(luan): Show a succes modal and refresh the buy list area
	c.Header("HX-Redirect", "/buy-list")
	c.Status(http.StatusOK)
}

func DeleteBuyList(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log(logger.INFO, "Error! %v", err)
		c.Status(http.StatusInternalServerError)
	}

	buyList := buylist.BuyList {
		ID: id,
	}
	err = buyList.Delete(db)
	if err != nil {
		logger.Log(logger.INFO, "Error! %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
