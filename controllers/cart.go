package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"xian-tao/models"
)

type CartController struct {
	beego.Controller
}

// 获取购物车数量
func GetCartGoodsCount(this *beego.Controller) int {
	userName := this.GetSession("userName")
	if userName == nil {
		return 0
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	conn, err := redis.Dial("tcp", ":6379")
	defer conn.Close()
	if err != nil {
		fmt.Println("redis连接失败")
		return 0
	}
	// 有多少种商品
	goodsSKUCount, _ := redis.Int(conn.Do("hlen", "cart_"+strconv.Itoa(user.Id)))
	return goodsSKUCount
}

// 处理添加购物车
func (this *CartController) HandleAddCart() {
	// 获取数据
	goodsSKUId, err1 := this.GetInt("goods_sku_id")
	goodsCount, err2 := this.GetInt("goods_count")
	userName := this.GetSession("userName")

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	// 校验数据
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["msg"] = "传递的数据不正确"
		this.Data["json"] = resp
		return
	}
	if userName == nil {
		resp["code"] = 2
		resp["msg"] = "当前用户未登录"
		this.Data["json"] = resp
		return
	}

	// 处理数据
	// 购物车数据存在 redis 中，用hash存储
	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	conn, err := redis.Dial("tcp", ":6379")
	defer conn.Close()
	if err != nil {
		resp["code"] = 3
		resp["msg"] = "Redis数据库连接错误"
		this.Data["json"] = resp
		return
	}
	// 先获取原来的数量，然后给数量加起来
	preGoodsCount, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), goodsSKUId))
	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), goodsSKUId, goodsCount+preGoodsCount)
	// 有多少种商品
	goodsSKUCount := GetCartGoodsCount(&this.Controller)
	resp["code"] = 0
	resp["msg"] = "ok"
	resp["goodsSKUCount"] = goodsSKUCount

	// 返回 json 数据
	this.Data["json"] = resp
}

// 显示 购物车页面
func (this *CartController) ShowCart() {
	userName := GetUser(&this.Controller)

	// 从redis种获取数据
	conn, err := redis.Dial("tcp", ":6379")
	defer conn.Close()
	if err != nil {
		fmt.Println("redis连接错误")
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")
	// 获取该用户购物车所有商品以及对应商品数量：map[string]int  map[商品id]商品数量
	cartGoods, _ := redis.IntMap(conn.Do("hgetall", "cart_"+strconv.Itoa(user.Id)))

	goodsSKUs := make([]map[string]interface{}, len(cartGoods))
	i := 0

	allGoodsPrice := 0 // 总金额
	allGoodsCount := 0 // 总件数

	for index, value := range cartGoods { // index: 商品id，value：商品数量
		goodsSKUId, _ := strconv.Atoi(index)
		var goodsSKU models.GoodsSKU
		goodsSKU.Id = goodsSKUId
		o.Read(&goodsSKU)

		goodsData := make(map[string]interface{})
		goodsData["goodsSKU"] = goodsSKU
		goodsData["goodsCount"] = value

		// 所有商品的总价
		allGoodsPrice += goodsSKU.Price * value
		// 所有商品件数
		allGoodsCount += value

		// 单个商品的总价（商品单价 * 商品数量）
		goodsData["goodsTotalPrice"] = goodsSKU.Price * value

		goodsSKUs[i] = goodsData
		i += 1
	}
	this.Data["allGoodsPrice"] = allGoodsPrice
	this.Data["allGoodsCount"] = allGoodsCount
	this.Data["goodsSKUs"] = goodsSKUs

	this.Data["title"] = "鲜淘驿站 - 购物车"
	this.TplName = "cart.html"
}

// 更新购物车数量
func (this *CartController) HandleUpdateCart() {
	// 获取数据
	goodsSKUId, err1 := this.GetInt("goods_sku_id")
	goodsCount, err2 := this.GetInt("goods_count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["msg"] = "获取数据失败"
		this.Data["json"] = resp
		return
	}
	userName := GetUser(&this.Controller)

	// 获取用户ID
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")

	conn, err := redis.Dial("tcp", ":6379")
	defer conn.Close()
	if err != nil {
		resp["code"] = 2
		resp["msg"] = "redis数据库连接失败"
		this.Data["json"] = resp
		return
	}
	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), goodsSKUId, goodsCount)
	resp["code"] = 0
	resp["msg"] = "ok"
	this.Data["json"] = resp
}

// 删除购物车商品
func (this *CartController) HandleDeleteCart() {
	goodsSKUId, err := this.GetInt("goods_sku_id")
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	if err != nil {
		resp["code"] = 1
		resp["msg"] = "获取商品信息失败"
		this.Data["json"] = resp
		return
	}

	userName := GetUser(&this.Controller)
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")

	conn, err := redis.Dial("tcp", ":6379")
	defer conn.Close()
	if err != nil {
		resp["code"] = 2
		resp["msg"] = "redis数据库连接失败"
		this.Data["json"] = resp
		return
	}
	conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), goodsSKUId)

	resp["code"] = 0
	resp["msg"] = "ok"

	this.Data["json"] = resp
}
