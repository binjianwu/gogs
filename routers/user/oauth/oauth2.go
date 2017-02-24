package oauth

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"golang.org/x/oauth2"
)

func HandleOauth2Login(conf *oauth2.Config, oauthStateString string, w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)

	log.Println("Redirect to github")
	log.Println("url:" + url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//According to the authentication server return code to get AccessToken, and returns oauth2Client
func HandleOauth2Callback(conf *oauth2.Config, oauthStateString string, w http.ResponseWriter, r *http.Request) *http.Client {
	ct := context.Background()
	state := r.FormValue("state")
	//oauthStateString is for oauth2 API calls to protect against CSRF
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	code := r.FormValue("code")
	fmt.Println("code:" + code)
	token, err := conf.Exchange(ct, code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("AccessToken:" + token.AccessToken)
	client := conf.Client(ct, token)
	return client
}
