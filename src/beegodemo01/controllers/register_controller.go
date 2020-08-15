package controllers

import (
	"beegodemo01/models"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io"
)

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) ShowRegister() {
	this.TplName = "register.html"
}

func (this *RegisterController) HandleRegister() {
	data := this.Ctx.Input.RequestBody
	registerUser := models.RegisterUser{}
	json.Unmarshal(data, &registerUser)
	resp := make(map[string]interface{})
	if registerUser.Username == "" || registerUser.Password == "" {
		resp["msg_no"] = 1
		resp["msg_content"] = "用户名密码不能为空"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	user := models.User{}
	user.Username = registerUser.Username

	o := orm.NewOrm()
	err := o.Read(&user, "username")
	if err == nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "用户名已存在"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	} else {
		w := md5.New()
		io.WriteString(w, registerUser.Password)
		password := fmt.Sprintf("%x", w.Sum(nil))
		user.Password = password
		_, err = o.Insert(&user)
		if err != nil {
			resp["msg_no"] = 1
			resp["msg_content"] = "创建用户失败"
			this.Data["json"] = resp
			this.ServeJSON()
			return
		} else {
			resp["msg_no"] = 2
			resp["msg_content"] = "创建用户成功"
			this.Data["json"] = resp
			this.ServeJSON()
			return
		}
	}
}
