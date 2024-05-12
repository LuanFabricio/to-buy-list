package to_buy_list

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tbl-backend/models/item"
	"tbl-backend/services/email"
)

func SendToBuyListEmail(db *sql.DB, buy_list_id int) error {
	to_buy_list, err := FetchToBuyList(db, buy_list_id)

	if err != nil {
		return err
	}
	if len(to_buy_list) <= 0 {
		return errors.New("Empty to buy list")
	}

	emailContent := `
	<h1>To Buy items</h1>
	<ol>
	`
	for i, to_buy_item := range to_buy_list {
		delta := int32(to_buy_item.CurrentQuantity) - int32(to_buy_item.MinQuantity)
		emailContent += fmt.Sprintf(
			"<li>%s: %d/%d(%d)</li>",
			to_buy_item.Name,
			to_buy_item.CurrentQuantity,
			to_buy_item.MinQuantity,
			delta,
		)

		if i+1 < len(to_buy_list) { emailContent += "\r\n" }
	}
	emailContent += "</ol>"

	log.Println(emailContent)

	to, err := email.FetchUsersEmail(db, buy_list_id)

	if err != nil {
		return err
	}

	return email.SendEmail(
		to,
		"To Buy List",
		emailContent+"IDK",
	)
}

func FetchToBuyList(db *sql.DB, buy_list_id int) ([]item.BuyItem, error) {
	buy_list_row, err := db.Query(`
	SELECT
	id, name, current_quantity, min_quantity, send_email
	FROM items
	WHERE buy_list_id = $1
	AND current_quantity >= min_quantity`,
	buy_list_id)
	defer buy_list_row.Close()

	if err != nil {
		return nil, err
	}

	var buy_list []item.BuyItem = []item.BuyItem{}
	var buy_item item.BuyItem
	for buy_list_row.Next() {
		buy_list_row.Scan(&buy_item.ID, &buy_item.Name, &buy_item.CurrentQuantity, &buy_item.MinQuantity, &buy_item.SendEmail)
		buy_list = append(buy_list, buy_item)
	}

	return buy_list, nil
}
