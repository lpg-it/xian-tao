package controllers

import "github.com/astaxie/beego"

type GoodsController struct {
	beego.Controller
}
// 显示主页
func (c *GoodsController) ShowIndex(){
	GetUser(&c.Controller)

	// 返回视图
	c.TplName = "index.html"
}