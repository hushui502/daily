package routers

import (
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"honeypot/manager/logger"
	"honeypot/manager/models"
)

func ListUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		users, _ := models.ListUser()
		ctx.Data["users"] = users
		ctx.Data["user"] = sess.Get("user")
		//log.Println(users)
		ctx.HTML(200, "user")
	} else {
		ctx.Redirect("/login/")
	}
}

func NewUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["user"] = sess.Get("user")
		ctx.HTML(200, "user_new")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoNewUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		// ctx.Req.ParseForm()
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		_ = models.NewUser(username, password)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}

func EditUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		Id := ctx.Params(":id")
		user, _ := models.GetUserById(Id)
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["user"] = user
		ctx.Data["username"] = user.UserName
		ctx.HTML(200, "user_edit")
	} else {
		ctx.Redirect("/login/")
	}
}

func DoEditUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		Id := ctx.Params(":id")
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		_ = models.UpdateUser(Id, username, password)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}

func DelUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		Id := ctx.Params(":id")
		_ = models.DelUser(Id)
		ctx.Redirect("/admin/user/")
	} else {
		ctx.Redirect("/login/")
	}
}

func Auth(ctx *macaron.Context, sess session.Store, cpt *captcha.Captcha) {
	if cpt.VerifyReq(ctx.Req) {
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		ret, err := models.Auth(username, password)
		logger.Logger.Println(ret, err)
		if err == nil && ret {
			_ = sess.Set("admin", username)
			ctx.Redirect("/admin/index/")
		} else {
			ctx.Redirect("/login/")
		}
	} else {
		message := "验证码输入错误"
		ctx.Data["message"] = message
		ctx.HTML(200, "error")
	}
}

func Logout(ctx *macaron.Context, sess session.Store) {
	_ = sess.Destory(ctx)
	ctx.Redirect("/login/")
}
