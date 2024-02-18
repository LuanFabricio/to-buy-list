package item

type BuyItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	CurrentQuantity uint32 `json:"current_quantity"`
	MinQuantity uint32 `json:"min_quantity"`
	SendEmail bool `json:"send_email"`
}
