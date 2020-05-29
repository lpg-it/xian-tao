package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"github.com/smartwalle/alipay"
	"strconv"
	"strings"
	"time"
	"xian-tao/models"
)

type OrderController struct {
	beego.Controller
}

// 展示订单页面
func (this *OrderController) ShowOrder() {
	// 获取数据
	userName := GetUser(&this.Controller)
	goodsSKUIds := this.GetStrings("goods_sku_id")

	// 校验数据
	if len(goodsSKUIds) == 0 {
		this.Redirect("/u/cart", 302)
		return
	}

	// 处理数据
	o := orm.NewOrm()
	// 获取用户数据
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")

	conn, _ := redis.Dial("tcp", ":6379")
	defer conn.Close()
	// 存储所有要结算的商品
	goodsBuffer := make([]map[string]interface{}, len(goodsSKUIds))

	allGoodsPrice := 0 // 结算总金额
	allGoodsCount := 0 // 总件数

	for index, value := range goodsSKUIds { // index：索引，value：结算商品ID(string)
		goodsData := make(map[string]interface{})
		goodsSKUId, _ := strconv.Atoi(value)
		// 查询商品数据
		var goodsSKU models.GoodsSKU
		goodsSKU.Id = goodsSKUId
		o.Read(&goodsSKU)

		goodsData["goodsSKU"] = goodsSKU
		// 获取商品数量
		goodsCount, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), value))
		goodsData["goodsCount"] = goodsCount
		// 计算小计：单个商品的总价（商品单价 * 商品数量）
		goodsTotalPrice := goodsSKU.Price * goodsCount
		goodsData["goodsTotalPrice"] = goodsTotalPrice

		// 计算总金额和件数
		allGoodsPrice += goodsTotalPrice
		allGoodsCount += goodsCount

		goodsBuffer[index] = goodsData
	}
	this.Data["goodsBuffer"] = goodsBuffer

	// 获取地址
	var addrs []models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Id", user.Id).All(&addrs)

	this.Data["addrs"] = addrs

	// 传递总金额和总件数
	this.Data["allGoodsPrice"] = allGoodsPrice
	this.Data["allGoodsCount"] = allGoodsCount

	freight := 10 // 运费
	this.Data["freight"] = freight
	this.Data["realPrice"] = allGoodsPrice + freight

	// 传递所有商品ID
	this.Data["goodsSKUIds"] = goodsSKUIds

	// 获取购物车数据
	goodsSKUCount := GetCartGoodsCount(&this.Controller)
	this.Data["goodsSKUCount"] = goodsSKUCount

	this.Data["title"] = "鲜淘驿站 - 订单"
	// 返回视图
	this.TplName = "place_order.html"
}

// 处理添加订单请求
func (this *OrderController) HandleAddOrder() {
	// 获取数据
	userName := GetUser(&this.Controller)

	addrId, _ := this.GetInt("addr_id")     // 地址ID
	payStyle, _ := this.GetInt("pay_style") // 支付方式

	goodsSKUIdsStr := this.GetString("goods_sku_ids") // 商品id（字符串）
	goodsSKUIds := strings.Split(goodsSKUIdsStr[1:len(goodsSKUIdsStr)-1], " ")

	totalGoodsCount, _ := this.GetInt("total_count")
	freight, _ := this.GetInt("freight")      // 运费
	realPrice, _ := this.GetInt("real_price") // 实付款

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	if len(goodsSKUIds) == 0 {
		resp["code"] = 1
		resp["userName"] = userName
		resp["msg"] = "未获取到商品"
		this.Data["json"] = resp
		return
	}

	// 处理数据
	o := orm.NewOrm()

	o.Begin() // 标记事务的开始

	// 向订单表中插入数据

	// 查询用户
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")
	// 查询地址
	var addr models.Address
	addr.Id = addrId
	o.Read(&addr)

	// 增加订单数据
	var orderInfo models.OrderInfo
	orderInfo.OrderId = strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(user.Id)
	orderInfo.User = &user
	orderInfo.Address = &addr
	orderInfo.PayMethod = payStyle
	orderInfo.TotalCount = totalGoodsCount
	orderInfo.TotalPrice = realPrice
	orderInfo.OrderStatus = 1
	orderInfo.TransitPrice = freight // 运费

	o.Insert(&orderInfo)

	// 向订单商品表中插入数据
	conn, _ := redis.Dial("tcp", ":6379")
	for _, value := range goodsSKUIds { // value: 商品id
		goodsSKUId, _ := strconv.Atoi(value)

		var goodsSKU models.GoodsSKU
		goodsSKU.Id = goodsSKUId

		i := 5
		for i > 0 {
			o.Read(&goodsSKU)

			var orderGoods models.OrderGoods

			orderGoods.OrderInfo = &orderInfo
			orderGoods.GoodsSKU = &goodsSKU

			goodsCount, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), goodsSKUId))

			if goodsCount > goodsSKU.Stock { // 库存不足
				resp["code"] = 2
				resp["userName"] = userName
				resp["msg"] = "商品：" + goodsSKU.Name + "库存不足"
				this.Data["json"] = resp
				o.Rollback() // 事务回滚
				return
			}
			preStock := goodsSKU.Stock

			orderGoods.Count = goodsCount
			orderGoods.Price = goodsCount * goodsSKU.Price
			o.Insert(&orderGoods)

			goodsSKU.Stock -= goodsCount
			goodsSKU.Sales += goodsCount

			// 返回更新成功的数量
			updateCount, _ := o.QueryTable("GoodsSKU").Filter("Id", goodsSKU.Id).Filter("Stock", preStock).Update(orm.Params{"Stock": goodsSKU.Stock, "Sales": goodsSKU.Sales})
			if updateCount == 0 {
				// 库存改变，更新失败
				if i > 0 {
					i -= 1
					continue
				}
				resp["code"] = 3
				resp["userName"] = userName
				resp["msg"] = "商品：" + goodsSKU.Name + "库存不足"
				this.Data["json"] = resp
				o.Rollback() // 事务回滚
				return
			} else {
				// 更新成功, 删除购物车对应商品
				conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), goodsSKUId)
				break
			}
		}

		o.Commit() // 提交事务
		// 返回数据
		resp["code"] = 0
		resp["userName"] = userName
		resp["msg"] = "ok"
		this.Data["json"] = resp
	}
}

// 处理支付
func (this *OrderController) HandlePay() {
	var aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlDCwH0XdbZEPv7LW2qB8/vh5Nk+bpzQtQTFh9IbROYrON4U6SFqj7HL3n5OGAmauYBWSEXzqUuAsGOKidT8vTQZeqArzLocqp00jIN+8GkXlXd/zuEcUhk2Fnz6opalROycXn5F61fVJf9AxalNMk8C+EVvUqxnP0CH1ih2ab8f/fl5slem1k+mTgwuesDimvaFXgAh91e7NjiqutG1oMpk6l2k+d/iddZh7egV4orn0vO0rYRM4qmSF2r6uLwCds0Hud7i1FFu6kBdQmYAjlNyFmixzP9QjPLqlz0gGxUqfnVjueqJQGN4kV6d8wf7N/SGLG6ryuyqukET0ExKoeQIDAQAB"
	var privateKey = "MIIEpAIBAAKCAQEAlDCwH0XdbZEPv7LW2qB8/vh5Nk+bpzQtQTFh9IbROYrON4U6SFqj7HL3n5OGAmauYBWSEXzqUuAsGOKidT8vTQZeqArzLocqp00jIN+8GkXlXd/zuEcUhk2Fnz6opalROycXn5F61fVJf9AxalNMk8C+EVvUqxnP0CH1ih2ab8f/fl5slem1k+mTgwuesDimvaFXgAh91e7NjiqutG1oMpk6l2k+d/iddZh7egV4orn0vO0rYRM4qmSF2r6uLwCds0Hud7i1FFu6kBdQmYAjlNyFmixzP9QjPLqlz0gGxUqfnVjueqJQGN4kV6d8wf7N/SGLG6ryuyqukET0ExKoeQIDAQABAoIBAGDksNPR16VDWxvJsIgEtZX1GzQyuyCJkil1Q4oh+H16T7mnp+MVOOdqiJRTXiUFxHBYykga+A+2Ob8PuI+W/7OKPav8dOLwSChZ3GUrRQ+csgs+WlocR8REveDQlG61FcLqnZyc/8cT+bnTg+v0iTZ2qRAqjhRN7T42ZhinoIoDJ8iSUSyTJ7leT34zjAy3XbMOc+UhfYdaPKtIGto3pOHQYINyoA8QVRi3vcPsbJAAfi85d+qUbBrwRoZuLUAcMe68pgB1JiiibFbk4TYoXtY1cYiTwTvBgV+TJLfqJ7gnjUgMHP6zBmOEANUOO5hrRa5/U/6fDb0SH5iBPzmQgfECgYEA2Srrw/W1UTeXRDpegT6VmQ7+pDRAsNFTP3lgsVQrjgZtGkQJy7t19U4HIIsmRgXLPETG4R3P/D8mgXtNJXTy/bX5u7Y2Nnj0qw5FiMRtbALtCLWFn9C/vxCuKRal8wBj3aHCTW0+qV2ruPWXQlfiRZ9LG1r88V3p0WGA6rzx+s0CgYEArrBAFrT0/4Tna5q+UbXjukAFw8Hk1NtpOOqBS4W/j3+YKuObXD1LLHDUi0WzGWxcrKdxfgpvIuCJhy9498d+E64Sn5Asiks5uPyr36zSTzKP1i7rywf5AT0vhij8BtzwnlfYbz8hV95jeLF1wXH5anF2eeWqZpcXH9O1Rp7svF0CgYEA1cmoIdiQb+zfED657Fg1I2GckwARsz/OyUzvQIMRAZcX7uSOFC9ul1gCMipqOkLX6XP3qYQUzUlJ2ewNbVNtJxDvUbi2M/ftPTwmfdaJtexHduxkKIlzSl/cY/y0z71Rksz8oAZsyoS5WbMD/j7QNSP053AyVFbUqNho9i2dtf0CgYBArZ89CQkBJmssyyGWTVsg1Z2Mylh4ezhtS15N4Rp4/gwQLS+TqloP/UKkwky6qAV0I5cAzMozRqGE/Q2z6BgFH1lj3NSw64NWu67DZVCE5DqfWcYR6UTHsajL6pbNz7YDWpEXN2+YAg4gXMw1sIZhY9sy7Nb3nw9/yDoBCMysPQKBgQDIFt2+HIbNAFu/mVcI4T3g1llDQtwiq87//tHOPt3/tSDF9y9dpdm7NXNaBao5mwzNtHxpg2h/OTthdypthARXTXc7QoAgL7DdVzrcXoaJbcq+COHTK7tOBQnurnsoflU4UltpwmWobjB9Yap8oMv6/yxISAlViMkR0kqyOGTLDw=="
	appId := "2016102400748415"
	client := alipay.New(appId, aliPublicKey, privateKey, false)

	// 获取数据
	orderId := this.GetString("order-id")
	totalPrice := this.GetString("total-price")

	var p = alipay.AliPayTradePagePay{}
	p.ReturnURL = "http://192.168.0.106:8080/u/pay-ok"
	p.Subject = "鲜淘驿站购物平台"
	p.OutTradeNo = orderId     // 订单编号
	p.TotalAmount = totalPrice // 总价
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payURL = url.String()
	this.Redirect(payURL, 302)
}

// 支付成功
func (this *OrderController) HandlePayOk() {
	// 获取数据
	orderId := this.GetString("out_trade_no")

	// 校验数据
	if orderId == "" {
		this.Redirect("/u/user-order", 302)
		return
	}
	// 处理数据
	o := orm.NewOrm()
	updateCount, _ := o.QueryTable("OrderInfo").Filter("OrderId", orderId).Update(orm.Params{"OrderStatus": 2})
	if updateCount == 0 {
		this.Redirect("/u/user-order", 302)
		return
	}

	// 返回视图
	this.Redirect("/u/user-order", 302)
}
