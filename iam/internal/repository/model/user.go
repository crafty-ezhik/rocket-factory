package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID      uuid.UUID
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UserRegistrationInfo struct {
	Info UserInfo
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type NotificationMethod struct {
	ProviderName string `db:"provider_name"`
	Target       string `db:"target"`
}
