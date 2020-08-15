package controllers

import (
	"beegodemo01/models"
	"encoding/json"
	"math"
	"strconv"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

type DeleteController struct {
	beego.Controller
}

func (this *DeleteController) DeleteArticle() {
	resp := make(map[string]interface{})
	data := this.Ctx.Input.RequestBody
	del := models.DelArticle{}
	err := json.Unmarshal(data, &del)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "转码失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	id, err := strconv.Atoi(del.Id)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "字符串转换为整型失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	article01 := models.Article{}
	article01.Id = id
	o := orm.NewOrm()
	_, err = o.Delete(&article01)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "删除失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	pageSize := 2
	qs := o.QueryTable("Article")
	var counts int64
	if del.ArticleType == "" {
		counts, err = qs.RelatedSel("ArticleType").Count()
		if err != nil {
			resp["msg_no"] = 1
			resp["msg_content"] = "查询文章数量失败"
			this.Data["json"] = resp
			this.ServeJSONP()
			return
		}
	} else {
		counts, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", del.ArticleType).Count()
		if err != nil {
			resp["msg_no"] = 1
			resp["msg_content"] = "查询文章数量失败"
			this.Data["json"] = resp
			this.ServeJSONP()
			return
		}
	}
	pageCount := math.Ceil(float64(counts) / float64(pageSize))
	nowpages := del.NowPage
	nowpage, err := strconv.Atoi(nowpages)
	var startIndex int
	if int(pageCount) >= nowpage {
		startIndex = (nowpage - 1) * pageSize
	} else {
		nowpage = nowpage - 1
		startIndex = (nowpage - 1) * pageSize
	}
	article02 := []models.Article{}
	if del.ArticleType == "" {
		qs.RelatedSel("ArticleType").Limit(pageSize, startIndex).All(&article02)
	} else {
		qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", del.ArticleType).Limit(pageSize, startIndex).All(&article02)
	}
	resp["msg_no"] = 2
	resp["msg_content"] = "删除成功"
	resp["article"] = article02
	resp["nowpage"] = nowpage
	resp["counts"] = counts
	resp["pageCount"] = pageCount
	this.Data["json"] = resp
	this.ServeJSON()
	return
}

func (this *DeleteController)DeleteType(){
	resp:=make(map[string]interface{})
	data:=this.Ctx.Input.RequestBody
	del:=models.DelType{}
	err:=json.Unmarshal(data,&del)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="转码失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	ids:=del.Id
	id,err:=strconv.Atoi(ids)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="类型转换失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	o:=orm.NewOrm()
	qs:=o.QueryTable("ArticleType")
	_,err=qs.Filter("Id",id).Delete()
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="删除文章类型失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	resp["msg_no"]=2
	resp["msg_content"]="成功"
	this.Data["json"]=resp
	this.ServeJSON()
	return
}
