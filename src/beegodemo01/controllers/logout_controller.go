package controllers

import "github.com/astaxie/beego"

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) HandleLogout() {
	this.DelSession("username")
	this.Redirect("/login",302)

}
