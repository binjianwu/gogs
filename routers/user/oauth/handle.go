package oauth

import (
	"log"
	"net/http"

	ct2 "context"

	"errors"

	"fmt"

	o2 "github.com/gogits/gogs/modules/auth/oauth2"
	"github.com/gogits/gogs/modules/context"
	"github.com/gogits/gogs/modules/setting"
)

var statueString string

func HandleLogin(oauthType string, w http.ResponseWriter, r *http.Request) {
	if setting.OauthService == nil {
		log.Println("oauth2 service not enabled")
		return
	}
	log.Println("oauth2 service is enabled")
	config, ok := o2.SocialMap[oauthType]
	if !ok {
		log.Println(oauthType + " oauth2 service not enabled")
		return
	}
	statueString = oauthType

	log.Println("clientID:" + config.ClientID)
	HandleOauth2Login(config, statueString, w, r)
}

func HandleCallback(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("state")
	switch t {
	case "github":
		c, ok := o2.SocialMap["github"]
		if !ok {
			ctx.Handle(404, "oauth2", errors.New("not found github oauth type"))
			return
		}
		token := HandleOauth2Callback(c, statueString, w, r)
		ctx.Session.Set("Oauth2AccessToken", token.AccessToken)
		client := c.Client(ct2.Background(), token)
		uj, err := o2.GetGithubUser(client)
		if err != nil {
			ctx.Handle(500, "oauth2", errors.New("get github user fail"))
			return
		}
		ctx.Session.Set("authType", "github")
		HandleSignIn(ctx, uj.Login)
	case "baidu":
		c, ok := o2.SocialMap["baidu"]
		if !ok {
			ctx.Handle(404, "oauth2", errors.New("baidu oauth2 service not enabled"))
			return
		}
		token := HandleOauth2Callback(c, statueString, w, r)
		ctx.Session.Set("Oauth2AccessToken", token.AccessToken)
		uj, err := o2.GetBaiduUser(token.AccessToken)
		if err != nil {
			ctx.Handle(500, "oauth2", err)
			return
		}
		ctx.Session.Set("authType", "github")
		HandleSignIn(ctx, uj.Uname)
	default:
		log.Println("error oauth2 type")
		ctx.Handle(404, "oauth2", errors.New("not found oauth type"))
		return
	}
}

func HandleRevoke(ctx *context.Context) {
	authType := fmt.Sprint(ctx.Session.Get("authType"))
	log.Println("auth type:" + authType)
	accessToken := fmt.Sprint(ctx.Session.Get("Oauth2AccessToken"))
	log.Println("token:" + accessToken)
	var revokeURL string
	var req *http.Request
	switch authType {
	case "github":
		c, ok := o2.SocialMap["github"]
		if !ok {
			ctx.Handle(404, "oauth2", errors.New("not found github oauth type"))
			return
		}
		revokeURL = "https://github.com/applications/" + c.ClientID + "/tokens/" + accessToken
		var err error
		req, err = http.NewRequest("DELETE", revokeURL, nil)
		if err != nil {
			log.Fatalln(err)
			ctx.Handle(500, "oauth2", errors.New("revoke token fail"))
			return
		}
	case "baidu":
		revokeURL = "https://openapi.baidu.com/rest/2.0/passport/auth/expireSession?access_token=" + accessToken
		var err error
		req, err = http.NewRequest("GET", revokeURL, nil)
		if err != nil {
			log.Fatalln(err)
			ctx.Handle(500, "oauth2", errors.New("revoke token fail"))
			return
		}
	default:
		ctx.Handle(400, "oauth2", errors.New("unknown auth type"))
		return
	}
	client := &http.Client{}
	log.Println("-----------------start to delete request------------------")
	_, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		ctx.Handle(500, "oauth2", errors.New("revoke token fail"))
		return
	}
}
