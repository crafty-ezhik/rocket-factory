package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID      uuid.UUID
	Info      UserInfo
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type UserRegistrationInfo struct {
	Info     UserInfo
	Password string
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}
