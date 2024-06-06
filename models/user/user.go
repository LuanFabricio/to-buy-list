package user

import (
	"database/sql"
	"tbl-backend/models/buy_list"
	"tbl-backend/services/logger"
)

type User struct {
	ID string `json:"id"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserDTO struct {
	Email string `form:"email" json:"email" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"confirm-password" json:"confirm-password" binding:"required"`
}

func (ud* UserDTO) Insert(db* sql.DB) (*User, error) {
	var user User
	var err = db.QueryRow(
		`INSERT INTO users (username, password, email) VALUES($1, $2, $3)
		RETURNING id, username, password`,
		ud.Username, ud.Password, ud.Email,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FetchUserById(db* sql.DB, id string) (User, error){
	var u User
	err := db.QueryRow(
		`SELECT id, username, password FROM users
		WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Username, &u.Password)

	return u, err
}

func FetchUserByUsername(db* sql.DB, username string) (User, error){
	var u User
	err := db.QueryRow(
		`SELECT id, username, password FROM users
		WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.Password)

	return u, err
}

func (u* User) FindByEmail(db* sql.DB, email string) (error) {
	err := db.QueryRow(
		`SELECT id, username, password FROM users
		WHERE username = $1`,
		email,
	).Scan(&u.ID, &u.Username, &u.Password)

	if err != nil {
		return err
	}

	return nil
}

func (u* User) FetchBuyLists(db* sql.DB) []buylist.BuyList {
	buyListArr := make([]buylist.BuyList, 0)

	rows, err := db.Query(
		`SELECT bl.id, bl.name, bl.owner_user_id FROM buy_list_access bla
		JOIN buy_list bl ON bl.id = bla.buy_list_id
		WHERE bla.user_id = $1`,
		u.ID,
	)
	if err != nil {
		logger.Log(logger.ERROR, "%v", err)
		return buyListArr
	}

	var buyList buylist.BuyList
	for rows.Next() {
		err = rows.Scan(&buyList.ID, &buyList.Name, &buyList.OwnerUserID)
		if err != nil {
			logger.Log(logger.WARNING, "%v", err)
		} else {
			buyListArr = append(buyListArr, buyList)
		}
	}

	return buyListArr
}
