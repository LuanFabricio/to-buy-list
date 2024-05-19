package buylist

import (
	"database/sql"
	"log"
	"tbl-backend/models/item"
)

type BuyList struct {
	ID int
	Name string
	OwnerUserID int
}

func (bl* BuyList) FetchByID(db* sql.DB, id int) error {
	err := db.QueryRow(
		`SELECT id, name, owner_user_id
		FROM buy_list
		WHERE id = $1`,
		id,
	).Scan(&bl.ID, &bl.Name, &bl.OwnerUserID)

	if err != nil {
		return err
	}

	return nil
}

func (bl* BuyList) FetchItems(db* sql.DB) ([]item.BuyItem, error){
	buyItemArr := make([]item.BuyItem, 0)

	rows, err := db.Query(
		`SELECT id, name, current_quantity, min_quantity, send_email FROM items
		WHERE buy_list_id = $1
		ORDER BY (current_quantity-min_quantity)`,
		bl.ID,
	)
	if err != nil {
		return buyItemArr, err
	}

	var buyItem item.BuyItem
	for rows.Next() {
		err = rows.Scan(&buyItem.ID, &buyItem.Name, &buyItem.CurrentQuantity, &buyItem.MinQuantity, &buyItem.SendEmail)

		if err != nil {
			log.Printf("[ERROR]: %v\n", err)
		} else {
			buyItemArr = append(buyItemArr, buyItem)
		}
	}

	return buyItemArr, nil
}
