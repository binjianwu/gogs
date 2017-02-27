package oauth2

import (
	"github.com/fabric8io/gogs/modules/log"
	"github.com/gogits/gogs/modules/setting"
	"golang.org/x/oauth2"
)

var SocialMap = make(map[string]*oauth2.Config)
var AllOauthes = []string{"github", "baidu", "qq", "twitter", "weibo"}

func NewOauthService() {
	if !setting.Cfg.Section("oauth2").Key("ENABLED").MustBool() {
		return
	}
	setting.OauthService = &setting.Oauther{}
	setting.OauthService.OauthInfos = make(map[string]*setting.OauthInfo)

	// Load all OAuth config data.
	for _, name := range AllOauthes {
		setting.OauthService.OauthInfos[name] = &setting.OauthInfo{
			ClientId:     setting.Cfg.Section("oauth2." + name).Key("CLIENT_ID").MustString(""),
			ClientSecret: setting.Cfg.Section("oauth2." + name).Key("CLIENT_SECRET").MustString(""),
			Scopes:       setting.Cfg.Section("oauth2." + name).Key("SCOPES").Strings(","),
			AuthUrl:      setting.Cfg.Section("oauth2." + name).Key("AUTH_URL").MustString(""),
			TokenUrl:     setting.Cfg.Section("oauth2." + name).Key("TOKEN_URL").MustString(""),
		}
		var endpoint oauth2.Endpoint
		endpoint.AuthURL = setting.OauthService.OauthInfos[name].AuthUrl
		endpoint.TokenURL = setting.OauthService.OauthInfos[name].TokenUrl
		SocialMap[name] = &oauth2.Config{
			ClientID:     setting.OauthService.OauthInfos[name].ClientId,
			ClientSecret: setting.OauthService.OauthInfos[name].ClientSecret,
			Scopes:       setting.OauthService.OauthInfos[name].Scopes,
			Endpoint:     endpoint,
			RedirectURL:  setting.AppUrl + "user/oauth2/sign_in",
		}
	}

	enabledOauths := make([]string, 0, 10)
	//github
	if setting.Cfg.Section("oauth2.github").Key("ENABLED").MustBool() {
		setting.OauthService.GitHub = true
		enabledOauths = append(enabledOauths, "GitHub")
	}

	//baidu
	if setting.Cfg.Section("oauth2.baidu").Key("ENABLED").MustBool() {
		setting.OauthService.Baidu = true
		enabledOauths = append(enabledOauths, "Baidu")
	}

	//qq
	if setting.Cfg.Section("oauth2.qq").Key("ENABLED").MustBool() {
		setting.OauthService.Tencent = true
		enabledOauths = append(enabledOauths, "QQ")
	}

	//twitter
	if setting.Cfg.Section("oauth2.twitter").Key("ENABLED").MustBool() {
		setting.OauthService.Twitter = true
		enabledOauths = append(enabledOauths, "Twitter")
	}

	//Weibo
	if setting.Cfg.Section("oauth2.weibo").Key("ENABLED").MustBool() {
		setting.OauthService.Weibo = true
		enabledOauths = append(enabledOauths, "Weibo")
	}
	log.Info("Oauth Service Enabled %s", enabledOauths)
}
