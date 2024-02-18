package item

import "fmt"

type BuyItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	CurrentQuantity uint32 `json:"current_quantity"`
	MinQuantity uint32 `json:"min_quantity"`
	SendEmail bool `json:"send_email"`
}

func (i BuyItem) Print() {
	fmt.Printf("Item id: %v\n", i.ID)
}
