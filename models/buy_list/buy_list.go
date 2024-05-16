package buylist

import (
	"database/sql"
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
