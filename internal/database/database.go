package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"secure-service/internal/service"
)

var Db *sql.DB

func InitDB() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		service.GetEnv("DB_HOST", "localhost"),
		service.GetEnv("DB_PORT", "5432"),
		service.GetEnv("DB_USER", "postgres"),
		service.GetEnv("DB_PASSWORD", "postgres"),
		service.GetEnv("DB_NAME", "secure_service"),
	)

	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err := Db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	return nil
}

func CloseDB() {
	if Db != nil {
		Db.Close()
	}
}

func CreateUser(email, username, passwordHash string) (*service.User, error) {
	row := Db.QueryRow(safeInsertQuery(usersTable), email, username, passwordHash)
	user := &service.User{}
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	return user, nil
}

func GetUserByID(userID int) (*service.User, error) {
	user := &service.User{}
	row := Db.QueryRow(safeLookupQuery(usersTable), userID)
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	return user, nil
}

func GetUserByEmail(email string) (*service.User, error) {
	user := &service.User{
		Email: email,
	}
	row := Db.QueryRow(safeFilterQuery(usersTable, user, false), email)
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	return user, nil
}

func UserExistsByEmail(email string) (bool, error) {
	user := &service.User{Email: email}
	row := Db.QueryRow(safeFilterQuery(usersTable, user, true), email)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("scan failed: %w", err)
	}
	return exists, nil
}

func GetDB() *sql.DB {
	return Db
}
