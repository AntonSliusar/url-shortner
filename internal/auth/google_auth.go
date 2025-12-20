package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//CLientId & ClientSecret add to config
func GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: "149514136876-3kfcfekcbrionav2fmo2jlrv3ko4f6hk.apps.googleusercontent.com",
		ClientSecret: "",
		RedirectURL: "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
