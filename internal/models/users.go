package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"
	sqlite3lib "modernc.org/sqlite/lib"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password []byte
	Created  time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (username, email, password)
	VALUES (?, ?, ?)`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3lib.SQLITE_CONSTRAINT_UNIQUE && strings.Contains(err.Error(), "users.email") {
				return ErrDuplicateEmails
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := `SELECT id, password FROM users WHERE email = ?`
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Get(id int) (User, error) {
	var u User
	stmt := `SELECT id, username, email, created_at FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		}
		return User{}, err
	}
	return u, nil
}
