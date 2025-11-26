package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserUUID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}

type WhoamiResponse struct {
	Session Session
	User    User
}
