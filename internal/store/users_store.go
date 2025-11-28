package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainPasswordText string) error {
	// 12 is cost factor
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPasswordText), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	p.plainText = &plainPasswordText
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	fmt.Println("Comparing password hash:", p.hash , "with plaintext password:", plaintextPassword)
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, fmt.Errorf("failed to compare passwords: %w", err)
		}
	}
	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserWithPasswordByID(id int64) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
	UpdatePassword(userID int64, newPassword string) error
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query :=
		`INSERT INTO users (email, password_hash, firstname, lastname, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id;
	`
	err := pg.db.QueryRow(query, user.Email, user.PasswordHash.hash, user.FirstName, user.LastName).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresUserStore) GetUserByID(id int64) (*User, error) {
	user := &User{PasswordHash: password{}}
	query := `
	SELECT id, email, firstname, lastname, created_at, updated_at from users where id = $1;
	`
	err := pg.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	user := &User{PasswordHash: password{}}
	query := `
	SELECT id, email, password_hash, firstname, lastname, created_at, updated_at from users where email = $1;
	`
	err := pg.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash.hash, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// bu kerak emas (faqat password change da ishlatildi)
func (pg *PostgresUserStore) GetUserWithPasswordByID(id int64) (*User, error) {
	user := &User{PasswordHash: password{}}
	query := `
	SELECT id, email, password_hash, firstname, lastname, created_at, updated_at from users where id = $1;
	`
	err := pg.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.PasswordHash.hash, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	query := `
	UPDATE users SET email = $1, password_hash = $2, firstname = $3, lastname = $4, updated_at = NOW()
	WHERE id = $5;
	`
	result, err := pg.db.Exec(query, user.Email, user.PasswordHash.hash, user.FirstName, user.LastName, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}
	return nil
}

func (pg *PostgresUserStore) DeleteUser(id int64) error {
	query := `
	DELETE FROM users WHERE id = $1;
	`
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("article with ID %d not found", id)
	}
	return nil
}

func (pg *PostgresUserStore) UpdatePassword(userID int64, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	query := `
	UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2;
	`
	result, err := pg.db.Exec(query, hashedPassword, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}
	return nil
}
 
