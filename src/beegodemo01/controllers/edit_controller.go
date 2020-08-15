package controllers

import (
	"beegodemo01/models"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type EditController struct {
	beego.Controller
}

func (this *EditController)ShowEdit(){
	this.TplName="update.html"
}

func (this *EditController)ShowEditContent(){
	resp:=make(map[string]interface{})
	data:=this.Ctx.Input.RequestBody
	edit:=models.Edit{}
	err:=json.Unmarshal(data,&edit)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="转码失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	ids:=edit.Id
	idi,err:=strconv.Atoi(ids)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="类型转换失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	article:=models.Article{Id: idi}
	o:=orm.NewOrm()
	err=o.Read(&article)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="读取文章失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	resp["msg_no"]=2
	resp["msg_content"]="成功"
	resp["article"]=article
	this.Data["json"]=resp
	this.ServeJSON()
}
func(this *EditController)HandleEditContent(){
	resp:=make(map[string]interface{})
	data:=this.Ctx.Input.RequestBody
	edit:=models.Edit{}
	err:=json.Unmarshal(data,&edit)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="转码失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}

	ids:=edit.Id
	idi,err:=strconv.Atoi(ids)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="类型转换失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	article:=models.Article{Id: idi}
	o:=orm.NewOrm()
	err=o.Read(&article)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="读取文章失败"
		this.Data["json"]=resp
		this.ServeJSON()
		return
	}
	beego.Info(article)
	article.Title=edit.Title
	article.Content=edit.Content

	if edit.ImgName!="" && edit.ImgContent!="" {
		imgmsg := edit.ImgContent
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

		imgnames := edit.ImgName
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
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModeAppend)
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

		article.Img="."+path
	}
	beego.Info(article)
	_,err=o.Update(&article)
	if err != nil {
		resp["msg_no"]=1
		resp["msg_content"]="更新文件出错"
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
