package model

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	UUID      uuid.UUID
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ExpiresAt *time.Time
}

type WhoamiResponse struct {
	Session Session
	User    User
}
