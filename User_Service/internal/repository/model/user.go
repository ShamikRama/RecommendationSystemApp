package model

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"username"`
	Password string `db:"password_hash"`
}

const (
	UserTable = "users"
)
