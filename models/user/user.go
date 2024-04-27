package user

import (
	"database/sql"
)

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hash bool `json:"hash"`
}

func (u* UserDTO) Insert(db* sql.DB) (*User, error) {
	var user User
	var err = db.QueryRow(
		`INSERT INTO users (username, password) VALUES($1, $2)
		RETURNING id, username, password`,
		u.Username, u.Password,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
