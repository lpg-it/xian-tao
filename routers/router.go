package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"xian-tao/controllers"
)

func init() {
	beego.InsertFilter("/u/*", beego.BeforeExec, loginFilter)
	// 注册
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleRed")
	// 登录
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	// 激活用户
	beego.Router("/active", &controllers.UserController{}, "get:ActiveUser")
	// 主页
	beego.Router("/", &controllers.GoodsController{}, "get:ShowIndex")
	// 退出登录
	beego.Router("/u/logout", &controllers.UserController{}, "get:Logout")
	// 用户中心信息
	beego.Router("/u/user-info", &controllers.UserController{}, "get:ShowUserInfo")
	// 用户中心全部订单
	beego.Router("/u/user-order", &controllers.UserController{}, "get:ShowUserOrder")
	// 用户中心收货地址
	beego.Router("/u/user-address", &controllers.UserController{}, "get:ShowUserAddress;post:HandleUserAddress")
	// 商品详情
	beego.Router("/goods-detail", &controllers.GoodsController{}, "get:ShowGoodsDetail")
	// 商品列表
	beego.Router("/goods-list", &controllers.GoodsController{}, "get:ShowGoodsList")
	// 商品搜索
	beego.Router("/goods-search", &controllers.GoodsController{}, "post:HandleGoodsSearch")
	// 添加购物车
	beego.Router("/u/add-cart", &controllers.CartController{}, "post:HandleAddCart")
	// 购物车页面
	beego.Router("/u/cart", &controllers.CartController{}, "get:ShowCart")
	// 更新购物车数量
	beego.Router("/u/update-cart", &controllers.CartController{}, "post:HandleUpdateCart")
	// 删除购物车商品数量
	beego.Router("/u/delete-cart", &controllers.CartController{}, "post:HandleDeleteCart")
	// 显示订单页面
	beego.Router("/u/order", &controllers.OrderController{}, "post:ShowOrder")
	// 添加订单
	beego.Router("/u/add-order", &controllers.OrderController{}, "post:HandleAddOrder")
	// 处理支付
	beego.Router("/u/pay", &controllers.OrderController{}, "get:HandlePay")
	// 支付成功
	beego.Router("/u/pay-ok", &controllers.OrderController{}, "get:HandlePayOk")
}

var loginFilter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
