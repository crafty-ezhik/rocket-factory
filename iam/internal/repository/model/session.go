package model

type Session struct {
	UserUUID  string `redis:"user_uuid"`
	CreatedAt int64  `redis:"createdAt"`
	UpdatedAt *int64 `redis:"updatedAt,omitempty"`
	ExpiresAt int64  `redis:"expiresAt"`
}
