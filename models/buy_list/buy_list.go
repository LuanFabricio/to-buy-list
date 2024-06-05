package buylist

import (
	"database/sql"
	"tbl-backend/models/item"
	"tbl-backend/services/logger"
)

type BuyList struct {
	ID int `form:"id" json:"id" binding:"-"`
	Name string `form:"name" json:"name" binding:"required"`
	OwnerUserID int `form:"ower_user_id" json:"ower_user_id" binding:"-"`
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
			logger.Log(logger.ERROR, "%v\n", err)
		} else {
			buyItemArr = append(buyItemArr, buyItem)
		}
	}

	return buyItemArr, nil
}

func (bl* BuyList) UserHaveAccess(db* sql.DB, userId string) (bool, error) {
	var rowCount int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM buy_list_access
		WHERE buy_list_id = $1
			AND user_id = $2`,
		bl.ID,
		userId,
	).Scan(&rowCount)
	if err != nil {
		return false, err
	}

	return rowCount >= 1, nil
}

func (bl* BuyList) AddAccessTo(db* sql.DB, userId string) error {
	insertString := `
		INSERT INTO buy_list_access (buy_list_id, user_id)
			VALUES($1, $2)
	`
	err := db.QueryRow(insertString, bl.ID, userId).Err()
	if err != nil {
		return err
	}

	return nil
}

func (bl* BuyList) Insert(db* sql.DB) (*BuyList, error) {
	logger.Log(logger.INFO, "Owner User ID: %d", bl.OwnerUserID)

	logger.Log(logger.INFO,
		"INSERT INTO buy_list (name, owner_user_id) VALUES(%v, %v) RETURNING id",
		bl.Name,
		bl.OwnerUserID,
	)

	insertString := `
		INSERT INTO buy_list (name, owner_user_id)
			VALUES($1, $2)
			RETURNING id
	`
	err := db.QueryRow(insertString, bl.Name, bl.OwnerUserID).Scan(&bl.ID)
	if err != nil {
		return nil, err
	}

	return bl, nil
}

func (bl* BuyList) Delete(db *sql.DB) error {
	deleteString := `
		DELETE FROM buy_list
		WHERE id = $1
		RETURNING id, name, owner_user_id
	`

	err := db.QueryRow(deleteString, bl.ID).Scan(&bl.ID, &bl.Name, &bl.OwnerUserID)
	return err
}

func FetchBuyListFromId(db *sql.DB, ID int) (*BuyList, error) {
	var bl BuyList

	selectString := `
		SELECT
			id, name, owner_user_id
		FROM buy_list
		WHERE id = $1
	`
	err := db.QueryRow(selectString, ID).Scan(&bl.ID, &bl.Name, &bl.OwnerUserID)
	if err != nil {
		return nil, err
	}

	return &bl, nil
}
