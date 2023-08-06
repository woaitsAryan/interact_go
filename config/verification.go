package config

import "time"

const (
	VERIFICATION_EMAIL_SUBJECT       = "OTP For Verification | Interact"
	VERIFICATION_EMAIL_BODY          = "OTP: "
	VERIFICATION_OTP_EXPIRATION_TIME = 10 * time.Minute
)
