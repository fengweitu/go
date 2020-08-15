package main

import (
	_ "beegodemo01/models"
	_ "beegodemo01/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("ShowPrePage", HandlePrePage)
	beego.AddFuncMap("ShowNextPage", HandleNextPage)
	beego.Run()
}

func HandlePrePage(index int) int {
	index = index - 1
	return index
}

func HandleNextPage(index int) int {
	index = index + 1
	return index
}
