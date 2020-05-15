package routers

import (
	"xian-tao/controllers"
	"github.com/astaxie/beego"
)

func init() {
    // 注册
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleRed")
    // 登录
    beego.Router("/login", &controllers.UserController{}, "get:ShowLogin")
    // 激活用户
    beego.Router("/active", &controllers.UserController{}, "get:ActiveUser")
}
