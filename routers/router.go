package routers

import (
	"github.com/astaxie/beego/context"
	"xian-tao/controllers"
	"github.com/astaxie/beego"
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
	beego.Router("/u/info", &controllers.UserController{}, "get:ShowUserInfo")
	// 用户中心全部订单
	beego.Router("/u/order", &controllers.UserController{}, "get:ShowUserOrder")
	// 用户中心收货地址
	beego.Router("/u/address", &controllers.UserController{}, "get:ShowUserAddress;post:HandleUserAddress")
	// 商品详情
	beego.Router("/goods-detail", &controllers.GoodsController{}, "get:ShowGoodsDetail")
	// 商品列表
	beego.Router("/goods-list", &controllers.GoodsController{}, "get:ShowGoodsList")
	// 商品搜索
	beego.Router("/goods-search", &controllers.GoodsController{}, "post:HandleGoodsSearch")
}

var loginFilter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}

