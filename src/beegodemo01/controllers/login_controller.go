package controllers

import (
	"beegodemo01/models"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) ShowLogin() {
	this.TplName = "login.html"
}

func (this *LoginController) SendCookie() {
	resp := make(map[string]interface{})
	cookie_username := this.Ctx.GetCookie("username")
	beego.Info("cookie:", cookie_username)
	resp["username"] = cookie_username
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *LoginController) HandleLogin() {
	resp := make(map[string]interface{})
	data := this.Ctx.Input.RequestBody
	if data == nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "ajax传给后台的数据为空"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	loginUser := models.LoginUser{}
	err := json.Unmarshal(data, &loginUser)
	if err != nil {
		beego.Info(err)
		resp["msg_no"] = 1
		resp["msg_content"] = "将ajax传入的数据解码失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	if loginUser.Check == true {
		this.Ctx.SetCookie("username", loginUser.Username, time.Second*3600)
	}
	if loginUser.Check == false {
		this.Ctx.SetCookie("username", loginUser.Username, -1)
	}
	this.SetSession("username", loginUser.Username)

	user := models.User{
		Username: loginUser.Username,
	}

	o := orm.NewOrm()
	err = o.Read(&user, "username")
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "用户名不存在"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	w := md5.New()
	io.WriteString(w, loginUser.Password)
	password := fmt.Sprintf("%x", w.Sum(nil))
	if user.Password != password {
		resp["msg_no"] = 1
		resp["msg_content"] = "用户名或密码出错"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	resp["msg_no"] = 2
	resp["msg_content"] = "登入成功"
	this.Data["json"] = resp
	this.ServeJSON()
	return
}
