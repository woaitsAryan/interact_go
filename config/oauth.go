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
		RedirectURL:  "/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	GoogleOAuthState = ""
	VALID_DOMAINS    = []string{"vitstudent.ac.in", "vit.ac.in"}
	NAME_REGEX       = `[1,2][0-9][A-Za-z]{3}[0-9]{4}`
)

func InitializeOAuthGoogle() {
	GoogleOAuthConfig.ClientID = initializers.CONFIG.GOOGLE_CLIENT_ID
	GoogleOAuthConfig.ClientSecret = initializers.CONFIG.GOOGLE_CLIENT_SECRET
	GoogleOAuthConfig.RedirectURL = initializers.CONFIG.BACKEND_URL + GoogleOAuthConfig.RedirectURL
	GoogleOAuthState = initializers.CONFIG.GOOGLE_OAUTH_STATE
}
