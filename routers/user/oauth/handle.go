package oauth

import (
	"log"
	"net/http"

	ct2 "context"

	"errors"

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
			log.Println("github oauth2 service not enabled")
			ctx.Handle(404, "type", errors.New("not found oauth type"))
			return
		}
		token := HandleOauth2Callback(c, statueString, w, r)
		ctx.Session.Set("Oauth2AccessToken", token.AccessToken)
		client := c.Client(ct2.Background(), token)
		uj, err := o2.GetGithubUser(client)
		if err != nil {
			ctx.Handle(500, "github", err)
			return
		}
		//ctx.Session.Set("authType", "github")
		HandleSignIn(ctx, uj.Login)
	case "baidu":
		c, ok := o2.SocialMap["baidu"]
		if !ok {
			log.Println("baidu oauth2 service not enabled")
			ctx.Handle(404, "type", errors.New("baidu oauth2 service not enabled"))
			return
		}
		token := HandleOauth2Callback(c, statueString, w, r)
		ctx.Session.Set("Oauth2AccessToken", token.AccessToken)
		uj, err := o2.GetBaiduUser(token.AccessToken)
		if err != nil {
			ctx.Handle(500, "", err)
			return
		}
		log.Println("name:" + uj.Uname)
		//ctx.Session.Set("authType", "baidu")
		HandleSignIn(ctx, uj.Uname)
	default:
		log.Println("error oauth2 type")
		ctx.Handle(404, "not found oauth type", nil)
		return
	}
}
