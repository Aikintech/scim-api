package facades

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauth struct {
	config   *oauth2.Config
	provider string
}

type SocialiteUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
}

var (
	GOOGLE_CLIENT_ID     = os.Getenv("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
	GOOGLE_REDIRECT_URL  = "http://localhost:3000/auth/google/callback"
	GOOGLE_USERINFO_URL  = "https://www.googleapis.com/oauth2/v2/userinfo"
	GOOGLE_USER_SCOPES   = []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"}
)

func Socialite() *oauth {
	return &oauth{}
}

func (o *oauth) Driver(provider string) *oauth {
	if provider == "google" {
		o.config = &oauth2.Config{
			ClientID:     GOOGLE_CLIENT_ID,
			ClientSecret: GOOGLE_CLIENT_SECRET,
			RedirectURL:  GOOGLE_REDIRECT_URL,
			Scopes:       GOOGLE_USER_SCOPES,
			Endpoint:     google.Endpoint,
		}
	}

	o.provider = provider

	return o
}

func (o *oauth) UserFromToken(token string) (SocialiteUser, error) {
	if o.provider == "google" {
		return o.getGoogleUserDetails(token)
	}
	return SocialiteUser{}, nil
}

func (o *oauth) getGoogleUserDetails(token string) (user SocialiteUser, err error) {
	tokenInfo, err := o.config.TokenSource(context.Background(), &oauth2.Token{AccessToken: token}).Token()
	if err != nil {
		return user, err
	}

	client := o.config.Client(context.Background(), tokenInfo)
	response, err := client.Get(GOOGLE_USERINFO_URL)
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return user, err
	}

	fmt.Println(userInfo)

	return SocialiteUser{
		ID:        userInfo["id"].(string),
		FirstName: userInfo["given_name"].(string),
		LastName:  userInfo["family_name"].(string),
		Email:     userInfo["email"].(string),
		Name:      userInfo["name"].(string),
		Avatar:    userInfo["picture"].(string),
	}, nil
}
