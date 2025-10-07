package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"altoai_mvp/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo() (UserRepo, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	// Create users table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating users table: %v", err)
	}

	return &postgresRepo{db: db}, nil
}

func (r *postgresRepo) List() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, email, name, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *postgresRepo) Get(id string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *postgresRepo) Create(email, name string) (models.User, error) {
	now := time.Now().UTC()
	u := models.User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := r.db.Exec(
		"INSERT INTO users (id, email, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		u.ID, u.Email, u.Name, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *postgresRepo) Update(id string, email, name *string) (models.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.User{}, err
	}
	defer tx.Rollback()

	u, err := r.Get(id)
	if err != nil {
		return models.User{}, err
	}

	if email != nil {
		u.Email = *email
	}
	if name != nil {
		u.Name = *name
	}
	u.UpdatedAt = time.Now().UTC()

	_, err = tx.Exec(
		"UPDATE users SET email = $1, name = $2, updated_at = $3 WHERE id = $4",
		u.Email, u.Name, u.UpdatedAt, u.ID,
	)
	if err != nil {
		return models.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *postgresRepo) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepo) Close() error {
	return r.db.Close()
}
