package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type psqlUserStore struct {
	db *sql.DB
}

func NewPSQLUserStore(db *sql.DB) *psqlUserStore {
	return &psqlUserStore{
		db: db,
	}
}

func (store *psqlUserStore) CreateUser(user *User) error {
	/*cmd := "INSERT INTO public.users (email, password_hash) VALUES ($1, $2)"
	_, err := store.db.Exec(cmd, user.Email, user.HashedPassword)*/
	cmd := "INSERT INTO public.users (email, password_hash) VALUES ($1, $2) RETURNING user_id"
	var userID int
	err := store.db.QueryRow(cmd, user.Email, user.HashedPassword).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	user.ID = userID
	return err
}

func (store *psqlUserStore) GetUserByEmail(email string) (*User, error) {

	query := "SELECT user_id, email, password_hash, created_at, updated_at FROM public.users WHERE email=$1"
	//row := store.db.QueryRow(query, email)
	//db := &sqlx.DB{}
	//db.QueryRow()

	//sqlscan.Select(ctx, db, &users, `SELECT id, name, email, age FROM users`)

	user := &User{}
	//sqlscan.ScanRow()

	err := sqlscan.Get(context.Background(), store.db, user, query, email)
	if err != nil {
		/*if err == sql.ErrNoRows {
			// Handle no rows case
			return nil, nil
		}*/
		return nil, err

	}
	return user, nil
	//sqlscan.ScanRow()
	//sqlscan.ScanOne(user, row)
	//row.Scan(user.)
	/*defer rows.Close()

	projects := []*Project{}
	for rows.Next() {
		var p Project
		err = rows.Scan(&p.Name, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil*/

}
