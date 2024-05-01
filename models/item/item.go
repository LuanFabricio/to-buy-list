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

func FindItem(id string, db* sql.DB) (*BuyItem, error) {
	row, err := db.Query("SELECT id, name, current_quantity, min_quantity, send_email FROM items WHERE id=$1", id)
	defer row.Close()

	if err != nil {
		return nil, err
	}

	var item BuyItem
	row.Next()
	err = row.Scan(&item.ID, &item.Name, &item.CurrentQuantity, &item.MinQuantity, &item.SendEmail)
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

func (bi* BuyItem) Update(db* sql.DB) (*BuyItem, error) {
	res, err := db.Query(`
		UPDATE items
		SET name = $1, current_quantity = $2, min_quantity = $3, send_email = $4
		WHERE id = $5
		RETURNING *`,
		bi.Name, bi.CurrentQuantity,
		bi.MinQuantity, bi.SendEmail,
		bi.ID,
	);

	if err != nil{
		return nil, err;
	}

	if res.Next() {
		res.Scan(&bi.ID, &bi.Name, &bi.CurrentQuantity, &bi.MinQuantity, &bi.SendEmail)
		return bi, nil
	}

	return nil, nil
}
