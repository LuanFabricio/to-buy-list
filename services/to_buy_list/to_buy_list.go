package to_buy_list

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tbl-backend/models/item"
	"tbl-backend/services/email"
)

func SendToBuyListToEveryone(db *sql.DB) {
	BUY_LISTS_QUERY := "SELECT id FROM buy_list"
	buyListRow, err := db.Query(BUY_LISTS_QUERY)
	defer buyListRow.Close()

	if err != nil {
		log.Fatal(err)
		return;
	}

	var buyListId int
	for buyListRow.Next() {
		err = buyListRow.Scan(&buyListId)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		} else {
			log.Printf("[INFO] Sending to buy list id %v", buyListId)
			err = SendToBuyListEmail(db, buyListId)
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
		}
	}
}

func SendToBuyListEmail(db *sql.DB, buyListId int) error {
	toBuyList, err := FetchToBuyList(db, buyListId)

	if err != nil {
		return err
	}
	if len(toBuyList) <= 0 {
		return errors.New("Empty to buy list")
	}

	emailContent := `
	<h1>To Buy items</h1>
	<ol>
	`
	for _, toBuyItem := range toBuyList {
		delta := int32(toBuyItem.CurrentQuantity) - int32(toBuyItem.MinQuantity)
		emailContent += fmt.Sprintf(
			"<li>%s: %d/%d(%d)</li>",
			toBuyItem.Name,
			toBuyItem.CurrentQuantity,
			toBuyItem.MinQuantity,
			delta,
		)
	}
	emailContent += "</ol>"
	log.Println(emailContent)

	to, err := email.FetchUsersEmail(db, buyListId)
	if err != nil {
		return err
	}

	return email.SendEmail(
		to,
		"To Buy List",
		emailContent,
	)
}

func FetchToBuyList(db *sql.DB, buyListId int) ([]item.BuyItem, error) {
	buy_list_row, err := db.Query(`
	SELECT
	id, name, current_quantity, min_quantity, send_email
	FROM items
	WHERE buy_list_id = $1
	AND current_quantity <= min_quantity`,
	buyListId)
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
