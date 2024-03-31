package model

import "time"

type EndpointAccess struct {
	ID        int64
	Address   string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}
