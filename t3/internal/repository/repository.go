package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"weather-api-t3/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// --- Users ---

func (r *Repository) CreateUser(ctx context.Context, name, email string) (*models.User, error) {
	var user models.User
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at, updated_at, deleted_at`
	err := r.db.GetContext(ctx, &user, query, name, email)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (r *Repository) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users WHERE deleted_at IS NULL`
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int, name, email string) (*models.User, error) {
	var user models.User
	query := `UPDATE users SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 AND deleted_at IS NULL RETURNING id, name, email, created_at, updated_at, deleted_at`
	err := r.db.GetContext(ctx, &user, query, name, email, id)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &user, nil
}

func (r *Repository) SoftDeleteUser(ctx context.Context, id int) error {
	query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// --- Cities ---

func (r *Repository) AddCity(ctx context.Context, userID int, cityName string) error {
	query := `INSERT INTO user_cities (user_id, city_name) VALUES ($1, $2) ON CONFLICT (user_id, city_name) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, cityName)
	if err != nil {
		return fmt.Errorf("add city: %w", err)
	}
	return nil
}

func (r *Repository) ListCities(ctx context.Context, userID int) ([]string, error) {
	var cities []string
	query := `SELECT city_name FROM user_cities WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &cities, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}
	return cities, nil
}

func (r *Repository) DeleteCity(ctx context.Context, userID int, cityName string) error {
	query := `DELETE FROM user_cities WHERE user_id = $1 AND city_name = $2`
	_, err := r.db.ExecContext(ctx, query, userID, cityName)
	if err != nil {
		return fmt.Errorf("delete city: %w", err)
	}
	return nil
}

// --- Weather History ---

func (r *Repository) SaveWeatherHistory(ctx context.Context, userID int, city string, temp float64, desc string) error {
	query := `INSERT INTO weather_history (user_id, city, temperature, description) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, userID, city, temp, desc)
	if err != nil {
		return fmt.Errorf("save weather history: %w", err)
	}
	return nil
}

func (r *Repository) GetWeatherHistory(ctx context.Context, userID int, city string, limit int) ([]models.WeatherRecord, error) {
	var history []models.WeatherRecord
	query := `SELECT temperature, description, requested_at FROM weather_history WHERE user_id = $1`
	args := []interface{}{userID}

	if city != "" {
		query += ` AND city = $2`
		args = append(args, city)
	}

	query += ` ORDER BY requested_at DESC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}

	err := r.db.SelectContext(ctx, &history, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get weather history: %w", err)
	}
	return history, nil
}
