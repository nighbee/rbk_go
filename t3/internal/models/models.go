package models

import "time"

type User struct {
	ID        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UserCity struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	CityName  string    `json:"city_name" db:"city_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type WeatherRecord struct {
	ID          int       `json:"id,omitempty" db:"id"`
	UserID      int       `json:"user_id,omitempty" db:"user_id"`
	City        string    `json:"city" db:"city"`
	Temperature float64   `json:"temperature" db:"temperature"`
	Description string    `json:"description" db:"description"`
	RequestedAt time.Time `json:"requested_at" db:"requested_at"`
}

type WeatherHistoryResponse struct {
	UserID  int             `json:"user_id"`
	City    string          `json:"city"`
	History []WeatherRecord `json:"history"`
}
