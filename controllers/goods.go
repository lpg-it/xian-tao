package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"xian-tao/models"
)

type GoodsController struct {
	beego.Controller
}

// 获取全部商品类型
func GetGoodsType(this *beego.Controller){
	o := orm.NewOrm()
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes
}

// 显示主页
func (this *GoodsController) ShowIndex() {
	userName := GetUser(&this.Controller)
	this.Data["userName"] = userName

	o := orm.NewOrm()
	// 展示 商品类型 数据
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes

	// 展示 首页轮播商品 数据
	var indexGoodsBanners []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&indexGoodsBanners)
	this.Data["indexGoodsBanners"] = indexGoodsBanners

	// 展示 首页促销商品 数据
	var indexPromotionBanners []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&indexPromotionBanners)
	this.Data["indexPromotionBanners"] = indexPromotionBanners

	// 展示 分类商品 数据
	goods := make([]map[string]interface{}, len(goodsTypes)) // 存储所有类型的商品
	// 插入 商品类型 数据
	for index, value := range goodsTypes {
		// 获取对应 商品类型的首页展示商品
		temp := make(map[string]interface{})
		temp["type"] = value
		goods[index] = temp
	}

	// 展示 文字商品 以及 图片商品
	for _, value := range goods {
		var textGoods []models.IndexTypeGoodsBanner
		var imageGoods []models.IndexTypeGoodsBanner
		// 获取 文字商品 数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").OrderBy("Index").Filter("GoodsType", value["type"]).Filter("DisplayType", 0).Limit(4, 0).All(&textGoods)
		// 获取 图片商品 数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").OrderBy("Index").Filter("GoodsType", value["type"]).Filter("DisplayType", 1).Limit(4, 0).All(&imageGoods)
		value["textGoods"] = textGoods
		value["imageGoods"] = imageGoods
	}
	this.Data["goods"] = goods

	// 返回视图
	this.TplName = "index.html"
}

// 显示商品详情页
func (this *GoodsController) ShowGoodsDetail() {
	userName := GetUser(&this.Controller)
	this.Data["userName"] = userName

	goodsId, err := this.GetInt("id")
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	o := orm.NewOrm()
	// 获取 商品SKU 数据
	var goodsSKU models.GoodsSKU
	goodsSKU.Id = goodsId
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType", "Goods").Filter("Id", goodsId).One(&goodsSKU)
	this.Data["goodsSKU"] = goodsSKU

	// 显示同类型的最新上架的两个商品
	var newGoodsSKUs []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType", goodsSKU.GoodsType).OrderBy("Time").Limit(2, 0).All(&newGoodsSKUs)
	this.Data["newGoodsSKUs"] = newGoodsSKUs

	// 获取所有商品类型
	GetGoodsType(&this.Controller)

	this.TplName = "detail.html"
}
