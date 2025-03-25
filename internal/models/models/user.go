package models

import (
	"database/sql"
	"time"
)

// User is structure for a single user in database
type (
	User struct {
		ID         string
		Name       string
		Email      string
		Password   string
		IsVerified bool
		VerifiedAt sql.NullTime
		CreatedAt  time.Time
	}
)
