package controllers

import (
	"beegodemo01/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type ContentController struct {
	beego.Controller
}

func (this *ContentController)ShowContent(){
	this.TplName="content.html"
}

func(this *ContentController)HandleContent(){
	resp:=make(map[string]interface{})
	data:=this.Ctx.Input.RequestBody
	detail:=models.Detail{}
	err:=json.Unmarshal(data,&detail)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="转码失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	ids:=detail.Id
	id,err:=strconv.Atoi(ids)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="类型转换失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	article:=models.Article{}
	article.Id=id
	o:=orm.NewOrm()
	err=o.Read(&article)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="读取文章失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	article.Count+=1
	_,err=o.Update(&article,"Count")
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="更新数据失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}

	username:=this.GetSession("username")
	user:=models.User{}
	user.Username=username.(string)
	err=o.Read(&user,"username")
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="查询用户失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	m2m:=o.QueryM2M(&article,"User")
	_,err=m2m.Add(&user)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="插入user失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	users:=[]models.User{}
	qs:=o.QueryTable("User")
	_,err=qs.Filter("Article__Article__Id",article.Id).Distinct().All(&users)
	if err != nil {
		beego.Info("查询出错")
		resp["msg_no"]=1
		resp["msg_content"]="查询出错"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}

	resp["msg_no"]=2
	resp["msg_content"]="查询详情成功"
	resp["article"]=article
	resp["users"]=users
	this.Data["json"]=resp
	this.ServeJSON()
	return
}
