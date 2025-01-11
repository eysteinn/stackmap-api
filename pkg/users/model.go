package users

import "time"

type User struct {
	ID             int       `db:"user_id"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"password_hash" json:"password_hash"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type UserStore interface {
	//GetUserByName(name string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
}
