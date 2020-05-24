package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math"
	"strconv"
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

// 分页控制
func pageTool(pageCount, pageIndex int) []int {
	var pages []int
	if pageCount <= 5 {  // 总页码数小于5，全部都显示
		pages = make([]int, pageCount)
		for i, _ := range pages {
			pages[i] = i + 1
		}
	} else if pageIndex <= 3 {  // 总页码数大于5，但是当前页码位于前三页
		pages = []int{1, 2, 3, 4, 5}
	} else if pageIndex >= pageCount - 3 {  // 总页码数大于5，但是当前页码位于后三页
		pages = []int{pageCount - 4, pageCount - 3, pageCount - 2, pageCount - 1, pageCount}
	} else {
		pages = []int{pageIndex - 2, pageIndex - 1, pageIndex, pageIndex + 1, pageIndex + 2}
	}
	if len(pages) == 0 {
		pages = append(pages, 1)
	}
	return pages
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

	// 添加登录用户的历史浏览记录
	userName := this.GetSession("userName")
	if userName == nil {
		// 未登录
		this.Data["userName"] = ""
	} else {
		// 已登录
		this.Data["userName"] = userName.(string)

		var user models.User
		user.Name = userName.(string)
		o.Read(&user, "Name")

		// 添加浏览记录
		conn, err := redis.Dial("tcp", ":6379")
		defer conn.Close()
		if err != nil {
			fmt.Println("redis连接失败")
		}
		// 把以前相同商品的浏览记录删除
		conn.Do("lrem", "history_" + strconv.Itoa(user.Id), 0, goodsId)
		// 添加新的商品浏览记录
		conn.Do("lpush", "history_" + strconv.Itoa(user.Id), goodsId)
	}

	this.TplName = "detail.html"
}

// 显示 商品列表 页
func (this *GoodsController) ShowGoodsList() {
	GetUser(&this.Controller)
	goodsTypeId, err := this.GetInt("type-id")
	if err != nil {
		fmt.Println("获取商品类型失败")
	}
	this.Data["goodsTypeId"] = goodsTypeId

	o := orm.NewOrm()
	// 获取当前商品类型名称
	var goodsType models.GoodsType
	goodsType.Id = goodsTypeId
	o.Read(&goodsType)
	this.Data["goodsType"] = goodsType

	// 获取所有商品类型
	GetGoodsType(&this.Controller)

	// 新品推荐
	var newGoodsSKUs []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).OrderBy("Time").Limit(2, 0).All(&newGoodsSKUs)
	this.Data["newGoodsSKUs"] = newGoodsSKUs

	// 商品分页
	// 对应类型商品总数量
	goodsCount, _ := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).Count()
	pageSize := 5  // 每一页显示多少个商品
	pageCount := math.Ceil(float64(int(goodsCount) / pageSize))  // 总共多少页
	pageIndex, err := this.GetInt("page-index")
	if err != nil {
		pageIndex = 1
	}
	this.Data["pageIndex"] = pageIndex

	// 商品列表
	var goodsSKUs []models.GoodsSKU
	start := (pageIndex - 1) * pageSize

	// 商品排序
	sortType := this.GetString("sort")
	if sortType == "price" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).OrderBy("Price").Limit(pageSize, start).All(&goodsSKUs)
		this.Data["goodsSKUs"] = goodsSKUs
	} else if sortType == "sale" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).OrderBy("Sales").Limit(pageSize, start).All(&goodsSKUs)
		this.Data["goodsSKUs"] = goodsSKUs
	}else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).Limit(pageSize, start).All(&goodsSKUs)
		this.Data["goodsSKUs"] = goodsSKUs
	}
	this.Data["sortType"] = sortType

	// 显示的页码
	pages := pageTool(int(pageCount), pageIndex)
	this.Data["pages"] = pages

	// 上一页
	prePageIndex := pageIndex - 1
	if prePageIndex <= 1 {
		prePageIndex = 1
	}
	this.Data["prePageIndex"] = prePageIndex

	// 下一页
	nextPageIndex := pageIndex + 1
	if nextPageIndex >= int(pageCount) {
		nextPageIndex = int(pageCount)
	}
	this.Data["nextPageIndex"] = nextPageIndex

	this.TplName = "list.html"
}

// 处理 搜索商品 结果
func (this *GoodsController) HandleGoodsSearch(){
	GetUser(&this.Controller)

	goodsSearchName := this.GetString("goods_search_name")
	o := orm.NewOrm()
	// 获取所有商品类型
	GetGoodsType(&this.Controller)

	// 展示商品搜索结果数据
	var goodsSKUs []models.GoodsSKU
	if goodsSearchName == "" {
		o.QueryTable("GoodsSKU").All(&goodsSKUs)
		this.Data["goodsSKUs"] = goodsSKUs
		this.TplName = "search.html"
		return
	}
	o.QueryTable("GoodsSKU").Filter("Name__icontains", goodsSearchName).All(&goodsSKUs)
	this.Data["goodsSKUs"] = goodsSKUs

	// 新品推荐
	var newGoodsSKUs []models.GoodsSKU
	o.QueryTable("GoodsSKU").OrderBy("Time").Limit(2, 0).All(&newGoodsSKUs)
	this.Data["newGoodsSKUs"] = newGoodsSKUs

	this.TplName = "search.html"
}