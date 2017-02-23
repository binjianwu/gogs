package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gogits/gogs/modules/context"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "4d0cdb1b78f753f5b275",
		ClientSecret: "a902f92592ccaa945e57503cf7a9f0ad3d8c8f27",
		// select level of access you want scopes
		Scopes:   []string{"user"},
		Endpoint: githuboauth.Endpoint,
	}
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = "xyz"
)

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	HandleOauth2Login(oauthConf, w, r)
}

func HandleGitHubCallback(ctx *context.Context, w http.ResponseWriter, r *http.Request) {

	client := HandleOauth2Callback(oauthConf, oauthStateString, w, r)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var uj = make(map[string]interface{})
	if err = json.Unmarshal(body, &uj); err != nil {
		log.Fatal(err)
		return
	}
	username, _ := uj["login"].(string)
	//id, _ := uj["id"].(int64)
	fmt.Println("username:" + username)
	HandleSignIn(ctx, username)

}
