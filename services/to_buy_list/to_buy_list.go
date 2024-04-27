package to_buy_list

import (
	"database/sql"
	"tbl-backend/models/item"
)

func FetchToBuyList(db *sql.DB) ([]item.BuyItem, error) {
	buy_list_row, err := db.Query("SELECT * FROM items")

	if err != nil {
		return nil, err
	}

	var buy_list []item.BuyItem = []item.BuyItem{}
	var buy_item item.BuyItem
	for buy_list_row.Next(){
		buy_list_row.Scan(&buy_item.ID, &buy_item.Name, &buy_item.CurrentQuantity, &buy_item.MinQuantity, &buy_item.SendEmail)
		if buy_item.CurrentQuantity <= buy_item.MinQuantity {
			buy_list = append(buy_list, buy_item)
		}
	}

	return buy_list, nil
}
