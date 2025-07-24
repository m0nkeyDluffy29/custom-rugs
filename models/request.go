package models

import "time"

type RugRequest struct {
	ID         int       `json:"id,omitempty"`
	Name       string    `json:"name" binding:"required"`
	EMAIL      string    `json:"email" binding:"required,email"`
	Details    string    `json:"details" binding:"required"`
	STATUS     string    `json:"status"`
	CREATED_AT time.Time `json:"created_at,omitempty"`
}
