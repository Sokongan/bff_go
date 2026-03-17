package app_domain

import (
	"time"

	"github.com/google/uuid"
)

type AppRegistry struct {
	ID           uuid.UUID `json:"id"`
	DSN          string    `json:"dsn"`
	RedirectPath string    `json:"redirect_path"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
