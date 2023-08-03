package config

import (
	"time"
)

const (
	ACCESS_TOKEN_TTL  = 5 * time.Minute
	REFRESH_TOKEN_TTL = 7 * 24 * time.Hour
)
