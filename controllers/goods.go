package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"xian-tao/models"
)

type GoodsController struct {
	beego.Controller
}
// 显示主页
func (this *GoodsController) ShowIndex(){
	userName := GetUser(&this.Controller)
	this.Data["userName"] = userName

	// 获取 商品类型 数据
	o := orm.NewOrm()
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes

	// 获取 首页轮播商品 数据
	var indexGoodsBanners []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&indexGoodsBanners)
	this.Data["indexGoodsBanners"] = indexGoodsBanners

	// 获取 首页促销商品 数据
	var indexPromotionBanners []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&indexPromotionBanners)
	this.Data["indexPromotionBanners"] = indexPromotionBanners

	// 返回视图
	this.TplName = "index.html"
}