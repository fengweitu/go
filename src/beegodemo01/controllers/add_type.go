package controllers

import (
	"beegodemo01/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type AddTypeController struct {
	beego.Controller
}

func (this *AddTypeController) Get() {
	this.TplName = "addType.html"
}

func (this *AddTypeController) ShowArticleType() {
	var articleType []models.ArticleType
	resp := make(map[string]interface{})
	o := orm.NewOrm()
	qs := o.QueryTable("article_type")
	_, err := qs.All(&articleType)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "查询文章类型出错"
		this.Data["josn"] = resp
		this.ServeJSON()
		return
	}

	username:=this.GetSession("username")

	resp["msg_no"] = 2
	resp["msg_content"] = "查询文章成功"
	resp["articletype"] = articleType
	resp["username"]=username
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *AddTypeController) HandleArticleType() {
	resp := make(map[string]interface{})
	data := this.Ctx.Input.RequestBody
	if data == nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "从前端接收的数据为空"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	articleType := models.ArticleType{}
	err := json.Unmarshal(data, &articleType)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "解码前端数据出错"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	o := orm.NewOrm()
	_, err = o.Insert(&articleType)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "添加文章类型出错"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	resp["msg_no"] = 2
	resp["msg_content"] = "添加文章类型成功"
	this.Data["json"] = resp
	this.ServeJSON()
	return
}
