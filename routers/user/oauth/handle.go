package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gogits/gogs/modules/context"
	"github.com/gogits/gogs/modules/setting"

	"github.com/gogits/gogs/modules/auth/oauth2"
)

var (
	githubStatueString = "xyz"
)

func HandleLogin(oauthType string, w http.ResponseWriter, r *http.Request) {
	if setting.OauthService == nil {
		log.Println("oauth2 service not enabled")
		return
	}
	log.Println("oauth2 service is enabled")
	config, ok := oauth2.SocialMap[oauthType]
	if !ok {
		log.Println(oauthType + " oauth2 service not enabled")
		return
	}
	log.Println("clientID:" + config.ClientID)
	HandleOauth2Login(config, githubStatueString, w, r)
}

func HandleGitHubCallback(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	config, ok := oauth2.SocialMap["github"]
	if !ok {
		log.Println("oauth2 service not enabled")
		return
	}
	client := HandleOauth2Callback(config, githubStatueString, w, r)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	//var uj = make(map[string]interface{})
	var uj oauth2.GithubUser
	if err = json.Unmarshal(body, &uj); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("login:" + uj.Login + "  ,   name:" + uj.Name)
	HandleSignIn(ctx, uj.Login)

}
