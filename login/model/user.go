package model

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID      int64      `json:"id"`
	Name    string     `json:"name"`
	Email   string     `json:"email"`
	Salt    string     `json:"salt"`
	Salted  string     `json:"salted"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
}

func (u *User) Insert(tx *sql.Tx, password string) (sql.Result, error) {
	stmt, err := tx.Prepare(`
insert into users (name, email, salt, salted)
values (?, ?, ?, ?)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	salt := Salt(100)

	return stmt.Exec(u.Name, u.Email, salt, Stretch(password, salt))
}

// Auth makes user authentication.
func Auth(db *sql.DB, email, password string) (*User, error) {
	rows, err := db.Query(`select id, name, email, salt, salted from users where email = ?`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u User
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Salt, &u.Salted); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if u.Salted != Stretch(password, u.Salt) {
		return nil, errors.New("user not found")
	}

	return &u, nil
}
