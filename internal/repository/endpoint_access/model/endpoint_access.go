package model

import "time"

// EndpointAccess ..
type EndpointAccess struct {
	ID        int64
	Address   string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
