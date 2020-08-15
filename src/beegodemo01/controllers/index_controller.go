package controllers

import (
	"beegodemo01/models"
	"encoding/json"
	"math"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get() {
	this.TplName = "index.html"
}

func (this *IndexController) ShowIndex() {

	//获取文章信息
	resp := make(map[string]interface{})
	o1 := orm.NewOrm()
	article := []models.Article{}
	qs1 := o1.QueryTable("article")

	//分页
	//获取文章总数
	counts, err := qs1.RelatedSel("ArticleType").Count()
	if err != nil {
		beego.Info("获取文章总数失败")
		resp["msg_no"] = 1
		resp["msg_content"] = "获取文章总数失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	//页数
	pageSize := 2
	pageCount := float64(counts) / float64(pageSize)
	pageCount = math.Ceil(pageCount)

	//设置首页文章
	_, err = qs1.RelatedSel("ArticleType").Limit(pageSize, 0).All(&article)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "查询文章失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	articleType := []models.ArticleType{}
	o2 := orm.NewOrm()
	qs2 := o2.QueryTable("article_type")
	_, err = qs2.All(&articleType)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "从数据库获取文章类型失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	username:=this.GetSession("username")

	resp["article_count"] = counts
	resp["page_count"] = pageCount
	resp["article"] = article
	resp["article_type"] = articleType
	resp["username"]=username
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *IndexController) HandleSelect() {
	resp := make(map[string]interface{})
	data := this.Ctx.Input.RequestBody
	if data == nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "类型为空"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	articleType := models.ArticleType{}
	err := json.Unmarshal(data, &articleType)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "转码失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	beego.Info(articleType.TypeName)
	o := orm.NewOrm()
	article := []models.Article{}
	qs := o.QueryTable("Article")
	counts, err := qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", articleType.TypeName).Count()
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "查询文章数量出错"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	pageSize := 2
	pageCount := math.Ceil(float64(counts) / float64(pageSize))
	_, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", articleType.TypeName).Limit(pageSize, 0).All(&article)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "查询文章失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	resp["msg_no"] = 2
	resp["msg_content"] = "查询成功"
	resp["article"] = article
	resp["article_count"] = counts
	resp["page_count"] = pageCount
	this.Data["json"] = resp
	this.ServeJSON()
	return
}

func (this *IndexController) HandleIndex() {
	resp := make(map[string]interface{})
	data := this.Ctx.Input.RequestBody
	paging := models.Paging{}
	err := json.Unmarshal(data, &paging)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "转码失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	pageSize := 2
	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var counts int64
	if paging.ArticleType == "" {
		counts, err = qs.RelatedSel("ArticleType").Count()
		if err != nil {
			resp["msg_no"] = 1
			resp["msg_content"] = "查询文章数量失败"
			this.Data["json"] = resp
			this.ServeJSONP()
			return
		}
	} else {
		counts, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", paging.ArticleType).Count()
		if err != nil {
			resp["msg_no"] = 1
			resp["msg_content"] = "查询文章数量失败"
			this.Data["json"] = resp
			this.ServeJSONP()
			return
		}
	}
	pageNum := math.Ceil(float64(counts) / float64(pageSize))
	nowpages := paging.NowPage
	nowpage, err := strconv.Atoi(nowpages)
	var startIndex int
	if paging.Handle == "first" {
		nowpage = 1
		startIndex = 0
	}
	if paging.Handle == "prev" {
		if nowpage == 1 {
			startIndex = 0
		} else {
			nowpage = nowpage - 1
			startIndex = (nowpage - 1) * pageSize
		}
	}
	if paging.Handle == "next" {
		if nowpage == int(pageNum) {
			startIndex = (nowpage - 1) * pageSize
		} else {
			nowpage = nowpage + 1
			startIndex = (nowpage - 1) * pageSize
		}
	}
	if paging.Handle == "last" {
		nowpage = int(pageNum)
		startIndex = (nowpage - 1) * pageSize
	}

	article := []models.Article{}
	if paging.ArticleType == "" {
		qs.RelatedSel("ArticleType").Limit(pageSize, startIndex).All(&article)
	} else {
		qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", paging.ArticleType).Limit(pageSize, startIndex).All(&article)
	}
	resp["msg_no"] = 2
	resp["msg_content"] = "成功"
	resp["article"] = article
	resp["nowpage"] = nowpage
	this.Data["json"] = resp
	this.ServeJSON()

}
