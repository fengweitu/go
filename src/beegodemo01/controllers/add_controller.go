package controllers

import (
	"beegodemo01/models"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"os"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

type AddArticleController struct {
	beego.Controller
}

func (this *AddArticleController) Get() {
	this.TplName = "add.html"
}

func (this *AddArticleController) ShowAddArticle() {
	resp := make(map[string]interface{})
	articleType := []models.ArticleType{}
	o := orm.NewOrm()
	qs := o.QueryTable("article_type")
	_, err := qs.All(&articleType)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "从数据库获取文章类型失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	username := this.GetSession("username")

	resp["msg_no"] = 2
	resp["msg_content"] = "从数据库获取文章类型成功"
	resp["articletype"] = articleType
	resp["username"] = username
	this.Data["json"] = resp
	this.ServeJSON()
	return
}

func (this *AddArticleController) HandleAddArticle() {
	//接受ajax传过来的数据
	resp := make(map[string]interface{})
	data := this.Ctx.Input.RequestBody
	if data == nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "接收数据为空"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	//将数据传给ArticleContent结构体
	articleContent := models.ArticleContnet{}
	err := json.Unmarshal(data, &articleContent)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "解码数据失败"
		this.Data["josn"] = resp
		this.ServeJSON()
		return
	}

	//处理图片
	imgmsg := articleContent.ImgContent
	imgdata := strings.Split(imgmsg, ",")
	imgBase64 := imgdata[1]
	img, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		beego.Info("base64 decode error:", err)
		resp["msg_no"] = 1
		resp["msg_content"] = "图片转码失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	imgnames := articleContent.ImgName
	imgname_suffix := path.Ext(imgnames)
	if imgname_suffix == "" {
		beego.Info("图片名出错")
		resp["msg_no"] = 1
		resp["msg_content"] = "图片名出错"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	sum := 0
	str := []string{".bmp", ".jpg", ".png", ".tif", ".gif", ".pcx", ".tga", ".exif", ".fpx", ".svg", ".psd", ".cdr", ".pcd", ".dxf", ".ufo", ".eps", ".ai", ".raw", ".WMF", ".webp"}
	for i := range str {
		if imgname_suffix == str[i] {
			sum++
		}
	}
	if sum != 1 {
		resp["msg_no"] = 1
		resp["msg_content"] = "文件类型出错"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	imgtime := time.Now().Format("2006-01-02 15-04-05")
	imgname := imgtime + imgname_suffix
	path := "./static/img/" + imgname
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModeAppend|os.ModePerm)
	if err != nil {
		beego.Info("创建文件失败:", err)
		resp["msg_no"] = 1
		resp["msg_content"] = "保存图片失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	_, err = w.WriteString(string(img))
	if err != nil {
		beego.Info("write error:", err)
		resp["msg_no"] = 1
		resp["msg_contnet"] = "写文件出错"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	err = w.Flush()
	if err != nil {
		beego.Info("刷新出错：", err)
		resp["msg_no"] = 1
		resp["msg_content"] = "刷新出错"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	articletype := models.ArticleType{TypeName: articleContent.ArticleType}
	o := orm.NewOrm()
	err = o.Read(&articletype, "TypeName")
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "文章类型不存在"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	article := models.Article{
		Title:       articleContent.ArticleTitle,
		Content:     articleContent.ArticleContent,
		Img:         "." + path,
		ArticleType: &articletype,
	}
	_, err = o.Insert(&article)
	if err != nil {
		resp["msg_no"] = 1
		resp["msg_content"] = "添加文章失败"
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	resp["msg_no"] = 2
	resp["msg_content"] = "成功"
	resp["data"] = article
	this.Data["json"] = resp
	this.ServeJSON()
	return
}
