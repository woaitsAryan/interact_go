package config

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  initializers.CONFIG.BACKEND_URL + "/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	GoogleOAuthState = ""
)

func InitializeOAuthGoogle() {
	GoogleOAuthConfig.ClientID = initializers.CONFIG.GOOGLE_CLIENT_ID
	GoogleOAuthConfig.ClientSecret = initializers.CONFIG.GOOGLE_CLIENT_SECRET
	GoogleOAuthState = initializers.CONFIG.GOOGLE_OAUTH_STATE
}
