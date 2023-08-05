package config

import "github.com/Pratham-Mishra04/interact/initializers"

var (
	MAILGUN_DOMAIN  = initializers.CONFIG.MAILGUN_DOMAIN
	MAILGUN_API_KEY = initializers.CONFIG.MAILGUN_PUBLIC_API_KEY
	SENDER          = "test"
	TEMPLATE_PATH   = ""
)
