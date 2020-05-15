package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"regexp"
	"strconv"
	"xian-tao/models"
)

type UserController struct {
	beego.Controller
}

// 显示注册页面
func (c *UserController) ShowReg(){
	c.TplName = "register.html"
}
// 处理注册数据
func (c *UserController) HandleRed(){
	// 获取数据
	userName := c.GetString("user_name")
	password := c.GetString("pwd")
	confirmPassword := c.GetString("cpwd")
	email := c.GetString("email")

	// 校验数据
	// 判断是否为空
	if userName == "" || password == "" || confirmPassword == "" || email == "" {
		c.Data["errMsg"] = "数据填写不完整，请重新输入~"
		c.TplName = "register.html"
		return
	}
	// 判断密码与确认密码是否一致
	if password != confirmPassword {
		c.Data["errMsg"] = "两次密码输入的不一致，请重新输入~"
		c.TplName = "register.html"
		return
	}
	// 使用正则判断邮箱格式
	regex, _ := regexp.Compile("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}")
	resEmail := regex.FindString(email)
	if resEmail == "" {
		c.Data["errMsg"] = "邮箱格式不正确，请重新输入~"
		c.TplName = "register.html"
		return
	}

	// 处理数据
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Password = password
	user.Email = email
	_, err := o.Insert(&user)
	if err != nil {
		c.Data["errMsg"] = "用户名已经存在，请重新输入~"
		c.TplName = "register.html"
		return
	}

	// 发送激活邮件



	emailConfig := `{"username":"431592976@qq.com","password":"nabvchjshgkacbbb","host":"smtp.qq.com","port":587}`
	emailConn := utils.NewEMail(emailConfig)
	emailConn.From = "431592976@qq.com"  // 指定发件人的邮箱地址
	emailConn.To = []string{user.Email}  // 指定收件人邮箱地址
	emailConn.Subject = "鲜淘驿站用户激活"  // 指定邮件的标题
	// 发给用户的是激活请求地址
	emailConn.HTML= user.Name + "您好，<br>"+"感谢您注册鲜淘驿站，为了避免您忘记账号或密码导致您的账户无法找回，请您验证Email地址。<br><a href=127.0.0.1:8080/active?id=" + strconv.Itoa(user.Id) + ">点击验证</a><br>------------<br>按钮无效？请复制粘贴下面的链接至浏览器地址栏打开：<br>"+"127.0.0.1:8080/active?id=" + strconv.Itoa(user.Id)
	err = emailConn.Send()
	if err != nil {
		c.Data["errMsg"] = "发送激活邮件失败，请重新发送~"
		c.TplName = "register.html"
		return
	}


	// 返回视图
	c.Ctx.WriteString("注册成功，请去邮箱激活用户。")


	// c.Redirect("/login", 302)



}
// 显示登录页面
func (c *UserController) ShowLogin(){
	c.TplName = "login.html"
}
// 激活用户
func (c *UserController) ActiveUser(){
	// 获取数据
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["errMsg"] = "要激活的用户不存在"
		c.TplName = "register.html"
		return
	}

	// 校验数据

	// 处理数据

	// 返回视图
}