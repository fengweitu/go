package routers

import (
	"beegodemo01/controllers"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
)

func init() {
	beego.InsertFilter("/v1/*", beego.BeforeRouter, FilterFunc)

	beego.Router("/", &controllers.MainController{})

	beego.Router("/register", &controllers.RegisterController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/sendcookie", &controllers.LoginController{}, "post:SendCookie")

	beego.Router("/v1/index", &controllers.IndexController{})
	beego.Router("/v1/index/content", &controllers.IndexController{}, "post:ShowIndex")
	beego.Router("/v1/index/paging", &controllers.IndexController{}, "post:HandleIndex")
	beego.Router("/v1/index/select", &controllers.IndexController{}, "post:HandleSelect")

	beego.Router("/v1/addtype", &controllers.AddTypeController{})
	beego.Router("/v1/addtype/content", &controllers.AddTypeController{}, "get:ShowArticleType;post:HandleArticleType")
	beego.Router("/v1/deltype", &controllers.DeleteController{}, "post:DeleteType")

	beego.Router("/v1/addarticle", &controllers.AddArticleController{})
	beego.Router("/v1/addarticle/content", &controllers.AddArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	// beego.Router("/v1/addarticle/img", &controllers.AddArticleController{}, "post:HandleImg")

	beego.Router("/v1/content", &controllers.ContentController{}, "get:ShowContent;post:HandleContent")

	beego.Router("/v1/edit", &controllers.EditController{}, "get:ShowEdit;post:ShowEditContent")
	beego.Router("/v1/edit/content", &controllers.EditController{}, "post:HandleEditContent")

	beego.Router("/v1/del", &controllers.DeleteController{}, "post:DeleteArticle")

	beego.Router("/logout", &controllers.LogoutController{}, "get:HandleLogout")
}

var FilterFunc = func(ctx *context.Context) {
	username := ctx.Input.Session("username")
	if username == nil {
		ctx.Redirect(302, "/login")

	}
}
