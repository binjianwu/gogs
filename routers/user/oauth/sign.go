package oauth

import (
	"log"
	"net/url"

	"github.com/gogits/gogs/modules/context"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/setting"
)

func HandleSignIn(ctx *context.Context, uname string) {
	user := &models.User{
		Name:      uname,
		Email:     uname + "@hand-china.com",
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

	u, err := models.UserSignInViaUname(uname)
	if err != nil {
		log.Println("user sign error")
		return
	}
	ctx.Session.Set("uid", u.ID)
	ctx.Session.Set("uname", u.Name)
	ctx.SetCookie(setting.CSRFCookieName, "", -1, setting.AppSubUrl)

	if redirectTo, _ := url.QueryUnescape(ctx.GetCookie("redirect_to")); len(redirectTo) > 0 {
		ctx.SetCookie("redirect_to", "", -1, setting.AppSubUrl)
		ctx.Redirect(redirectTo)
		return
	}
	ctx.Redirect(setting.AppSubUrl + "/")
}
