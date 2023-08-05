package config

import "time"

const (
	EMAIL_SUBJECT       = "OTP For Verification | Interact"
	EMAIL_BODY          = "OTP: "
	OTP_EXPIRATION_TIME = 10 * time.Minute
)
