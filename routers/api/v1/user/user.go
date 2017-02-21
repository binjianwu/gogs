// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/Unknwon/com"

	api "github.com/gogits/go-gogs-client"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/auth"
	"github.com/gogits/gogs/modules/base"
	"github.com/gogits/gogs/modules/context"
)

func Search(ctx *context.APIContext) {
	opts := &models.SearchUserOptions{
		Keyword:  ctx.Query("q"),
		Type:     models.USER_TYPE_INDIVIDUAL,
		PageSize: com.StrTo(ctx.Query("limit")).MustInt(),
	}
	if opts.PageSize == 0 {
		opts.PageSize = 10
	}

	users, _, err := models.SearchUserByName(opts)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	results := make([]*api.User, len(users))
	for i := range users {
		results[i] = &api.User{
			ID:        users[i].ID,
			UserName:  users[i].Name,
			AvatarUrl: users[i].AvatarLink(),
			FullName:  users[i].FullName,
		}
		if ctx.IsSigned {
			results[i].Email = users[i].Email
		}
	}

	ctx.JSON(200, map[string]interface{}{
		"ok":   true,
		"data": results,
	})
}

func GetInfo(ctx *context.APIContext) {
	u, err := models.GetUserByName(ctx.Params(":username"))
	if err != nil {
		if models.IsErrUserNotExist(err) {
			ctx.Status(404)
		} else {
			ctx.Error(500, "GetUserByName", err)
		}
		return
	}

	// Hide user e-mail when API caller isn't signed in.
	if !ctx.IsSigned {
		u.Email = ""
	}
	ctx.JSON(200, u.APIFormat())
}

func GetAuthenticatedUser(ctx *context.APIContext) {
	ctx.JSON(200, ctx.User.APIFormat())
}

//create access token by username and password
func CreateBaseTokens(ctx *context.APIContext, form auth.SignInForm) {
	_, err := models.UserSignIn(form.UserName, form.Password)
	if err != nil {
		ctx.JSON(401, "error:Invalid username or password")
		return
	}
	var t struct{ Authorization string }
	v := base.BasicAuthEncode(form.UserName, form.Password)
	t.Authorization = "Basic " + v
	ctx.JSON(201, t)

}
