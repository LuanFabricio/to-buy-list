package views

import "tbl-backend/models/item"

type ViewIndex struct {
	BuyItems []item.BuyItem
	ToBuyItems []item.BuyItem
}
