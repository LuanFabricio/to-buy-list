package item

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BuyItem struct {
	ID string `json:"id"`
	BuyListId string `json:"list_id"`
	Name string `json:"name"`
	CurrentQuantity uint32 `json:"current_quantity"`
	MinQuantity uint32 `json:"min_quantity"`
	SendEmail bool `json:"send_email"`
}

func FindItems(db* sql.DB, order_by_id bool)  ([]BuyItem, error) {
	select_query := "SELECT id, name, current_quantity, min_quantity, send_email, buy_list_id  FROM items"

	if order_by_id {
		select_query += " ORDER BY id"
	}

	row, err := db.Query(select_query)
	defer row.Close()

	if err != nil {
		return nil, err
	}

	var item BuyItem
	items := []BuyItem{}
	for row.Next() {
		err = row.Scan(
			&item.ID,
			&item.Name,
			&item.CurrentQuantity,
			&item.MinQuantity,
			&item.SendEmail,
			&item.BuyListId,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func FindItem(id string, db* sql.DB) (*BuyItem, error) {
	row, err := db.Query(`
		SELECT
			id, name,
			current_quantity, min_quantity,
			send_email, buy_list_id
		FROM items WHERE id=$1`, id)
	defer row.Close()

	if err != nil {
		return nil, err
	}

	var item BuyItem
	row.Next()
	err = row.Scan(
		&item.ID,
		&item.Name,
		&item.CurrentQuantity,
		&item.MinQuantity,
		&item.SendEmail,
		&item.BuyListId,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (bi* BuyItem) Insert(db* sql.DB) (*BuyItem, error) {
	var (
		id uint32
		name string
		current_quantity uint64
		min_quantity uint64
		send_email bool
		buy_list_id uint32
	)

	err := db.QueryRow(
		`INSERT INTO items (name, current_quantity, min_quantity, send_email, buy_list_id)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, name, current_quantity, min_quantity, send_email, buy_list_id`,
		bi.Name, bi.CurrentQuantity, bi.MinQuantity, bi.SendEmail, bi.BuyListId,
	).Scan(&id, &name, &current_quantity, &min_quantity, &send_email, &buy_list_id)

	if err != nil {
		return nil, err
	}

	return &BuyItem {
		ID: fmt.Sprint(id),
		Name: name,
		CurrentQuantity: uint32(current_quantity),
		MinQuantity: uint32(min_quantity),
		SendEmail: send_email,
		BuyListId: fmt.Sprint(buy_list_id),
	}, nil
}

func (bi* BuyItem) Update(db* sql.DB) (*BuyItem, error) {
	res, err := db.Query(`
		UPDATE items
		SET
			name = $1,
			current_quantity = $2,
			min_quantity = $3,
			send_email = $4,
			buy_list_id = $5
		WHERE id = $6
		RETURNING *`,
		bi.Name, bi.CurrentQuantity,
		bi.MinQuantity, bi.SendEmail,
		bi.BuyListId, bi.ID,
	);

	if err != nil{
		return nil, err;
	}

	if res.Next() {
		err := res.Scan(
			&bi.ID,
			&bi.Name,
			&bi.CurrentQuantity,
			&bi.MinQuantity,
			&bi.SendEmail,
			&bi.BuyListId,
		)
		if err != nil {
			return nil, err
		}
		return bi, nil
	}

	return nil, nil
}

func (bi* BuyItem) LoadFromForm(c *gin.Context) {
	if postId := c.PostForm("id"); postId != "" {
		bi.Name = postId
	}

	if postName := c.PostForm("name"); postName != "" {
		bi.Name = postName
	}

	if currentQuantity := c.PostForm("current_quantity"); currentQuantity != "" {
		currentQuantity, _ := strconv.ParseUint(currentQuantity, 10, 32)
		bi.CurrentQuantity = uint32(currentQuantity)
	}

	if minQuantity := c.PostForm("min_quantity"); minQuantity != "" {
		minQuantity, _ := strconv.ParseUint(minQuantity, 10, 32)
		bi.MinQuantity = uint32(minQuantity)
	}

	if buyListId := c.PostForm("buy_list_id"); buyListId != "" {
		bi.BuyListId = buyListId
	}

	postSendEmail := c.PostForm("send_email")
	bi.SendEmail = postSendEmail == "on"

}
