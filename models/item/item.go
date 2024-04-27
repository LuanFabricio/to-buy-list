package item

import (
	"database/sql"
	"fmt"
)

type BuyItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	CurrentQuantity uint32 `json:"current_quantity"`
	MinQuantity uint32 `json:"min_quantity"`
	SendEmail bool `json:"send_email"`
}

func FindItems(db* sql.DB)  ([]BuyItem, error) {

	row, err := db.Query("SELECT id, name, current_quantity, min_quantity, send_email FROM items")
	defer row.Close()

	if err != nil {
		return nil, err
	}

	var item BuyItem
	items := []BuyItem{}
	for row.Next() {
		err = row.Scan(&item.ID, &item.Name, &item.CurrentQuantity, &item.MinQuantity, &item.SendEmail)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (bi* BuyItem) Insert(db* sql.DB) (*BuyItem, error) {
	var (
		id uint32
		name string
		current_quantity uint64
		min_quantity uint64
		send_email bool
	)

	err := db.QueryRow(
		`INSERT INTO items (name, current_quantity, min_quantity, send_email)
		VALUES($1, $2, $3, $4)
		RETURNING id, name, current_quantity, min_quantity, send_email`,
		bi.Name, bi.CurrentQuantity, bi.MinQuantity, bi.SendEmail,
	).Scan(&id, &name, &current_quantity, &min_quantity, &send_email)

	if err != nil {
		return nil, err
	}

	return &BuyItem {
		ID: fmt.Sprint(id),
		Name: name,
		CurrentQuantity: uint32(current_quantity),
		MinQuantity: uint32(min_quantity),
		SendEmail: send_email,
	}, nil
}