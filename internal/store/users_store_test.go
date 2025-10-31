package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=articles sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	_, err = db.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}
	return db
}

func TestCreateUser(t *testing.T) {
	db := setupUserDB(t)
	defer db.Close()
	store := NewPostgresUserStore(db)
	tests := []struct {
		name    string
		user *User
		wantErr bool
	}{
		{
			name: "Valid User",
			user: &User{
				Email:       "c2KQw@example.com",
				PasswordHash: "hashed_password",
				FirstName:   "John",
				LastName:    "Doe",
			},
			wantErr: false,
		},
		{
			name: "Invalid User",
			user: &User{
				Email:       "invalid_email",
				PasswordHash: "hashed_password",
				FirstName:   "John",
				LastName:    "Doe",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdUser, err := store.CreateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err, "Expected error but got none")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.user.Email, createdUser.Email)
			assert.Equal(t, tt.user.PasswordHash, createdUser.PasswordHash)
			assert.Equal(t, tt.user.FirstName, createdUser.FirstName)
			assert.Equal(t, tt.user.LastName, createdUser.LastName)
			assert.NotZero(t, createdUser.ID)
		})
	}
}
