package item

import (
	"log"
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

	return []BuyItem{}, nil
}

func (bi* BuyItem) Insert(db* sql.DB) (*BuyItem, error) {
	var id uint32
	db.QueryRow(
		`INSERT INTO items (name, current_quantity, min_quantity, send_email)
		VALUES($1, $2, $3, $4)
		RETURNING id`,
		bi.Name, bi.CurrentQuantity, bi.MinQuantity, bi.SendEmail).Scan(&id)

	bi.ID = fmt.Sprint(id)

	row := db.QueryRow("SELECT * FROM items WHERE id = $1", id)

	var (
		name string
		current_quantity uint64
		min_quantity uint64
		send_email bool
	)
	if err := row.Scan(&id, &name, &current_quantity, &min_quantity, &send_email); err != nil {
		log.Fatal(err)
	}

	return &BuyItem {
		ID: fmt.Sprint(id),
		Name: name,
		CurrentQuantity: uint32(current_quantity),
		MinQuantity: uint32(min_quantity),
		SendEmail: send_email,
	}, nil
}
