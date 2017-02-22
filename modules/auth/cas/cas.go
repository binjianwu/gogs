package cas

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/context"
	"github.com/gogits/gogs/modules/setting"
	"gopkg.in/cas.v1"
)

var cas_server string = "https://login.hand-china.com/sso/"

//Redirect to cas login
func CasLogin(ctx *context.Context, w http.ResponseWriter, r *http.Request) {

	u, _ := url.Parse(cas_server)
	client := cas.NewClient(&cas.Options{
		URL: u,
	})

	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cas.IsAuthenticated(r) {
			client.RedirectToLogin(w, r)
			return
		}
		//获取用户名
		username := cas.Username(r)

		//user info
		user := &models.User{
			Name:      username,
			Email:     username + "@hand-china.com",
			Passwd:    "handhand",
			IsActive:  true,
			LoginType: models.LOGIN_NOTYPE,
		}
		//create new user
		if err := models.CreateUser(user); err != nil {
			log.Println(err)
		} else {
			log.Println("create user success")
		}

		//login
		if SignIn(ctx, username) {
			// Clear whatever CSRF has right now, force to generate a new one
			ctx.SetCookie(setting.CSRFCookieName, "", -1, setting.AppSubUrl)

			if redirectTo, _ := url.QueryUnescape(ctx.GetCookie("redirect_to")); len(redirectTo) > 0 {
				ctx.SetCookie("redirect_to", "", -1, setting.AppSubUrl)
				ctx.Redirect(redirectTo)
				return
			}
			ctx.Redirect(setting.AppSubUrl + "/")
		} else {
			client.RedirectToLogin(w, r)
		}

	})

	handler.ServeHTTP(w, r)

	//ctx.Session.Set("ticket", ticket)

	//根据ticker从cas服务器获取用户信息返回

	//判断该用户在gogs是否存在，不存在则创建新用户，密码默认。存在则调用gogs的登录

	//创建新用户后直接调用gogs的登录方法，生成本地session

}

//Redirect to cas logout
func CasLogout(ctx *context.Context, w http.ResponseWriter, r *http.Request) {

	u, _ := url.Parse(cas_server)
	client := cas.NewClient(&cas.Options{
		URL:         u,
		SendService: true,
	})
	referer := r.Referer()

	log.Println("from:" + referer)

	r.URL.Path = "/user/login"

	handler := client.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		client.RedirectToLogout(w, r)
		return
	})

	handler.ServeHTTP(w, r)

}

//gogs cas user sign ,only use username
func SignIn(ctx *context.Context, username string) bool {
	u, err := models.UserSignInViaCas(username)
	if err != nil {
		log.Println("user sign error")
		return false
	}
	ctx.Session.Set("uid", u.ID)
	ctx.Session.Set("uname", u.Name)
	return true

}
