package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"regexp"
	"strconv"
	"time"
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
	emailConn.HTML= `尊敬的` + user.Name + `，您好<br><br>感谢您注册鲜淘驿站，为了避免您忘记账号或密码导致您的账户无法找回，请您验证Email地址。<br><br>请复制粘贴下面的链接至浏览器地址栏打开：<br><br>127.0.0.1:8080/active?id=` + strconv.Itoa(user.Id) + `<br>`
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
// 激活用户
func (c *UserController) ActiveUser(){
	// 获取数据
	id, err := c.GetInt("id")

	// 校验数据
	if err != nil {
		c.Data["errMsg"] = "要激活的用户不存在"
		c.TplName = "register.html"
		return
	}
	// 处理数据
	o := orm.NewOrm()
	var user models.User
	user.Id = id
	err = o.Read(&user)
	if err != nil {
		c.Data["errMsg"] = "要激活的用户不存在"
		c.TplName = "register.html"
		return
	}
	user.Active = true
	o.Update(&user)
	// 返回视图
	c.Redirect("/login", 302)
}
// 显示登录页面
func (c *UserController) ShowLogin(){
	userName := c.Ctx.GetCookie("userName")
	// base64解密
	tempUserName, _ := base64.StdEncoding.DecodeString(userName)
	if string(tempUserName) == "" {
		// 没有记住用户名
		c.Data["userName"] = ""
		c.Data["checked"] = ""
	}else{
		c.Data["userName"] = string(tempUserName)
		c.Data["checked"] = "checked"
	}
	c.TplName = "login.html"
}
// 处理登录数据
func (c *UserController) HandleLogin(){
	// 获取数据
	userName := c.GetString("username")
	password := c.GetString("pwd")
	// 校验数据
	if userName == "" || password == ""{
		c.Data["errMsg"] = "用户名或密码不能为空"
		c.TplName = "login.html"
		return
	}

	// 处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		c.Data["errMsg"] = "用户名或密码错误"
		c.Data["userName"] = userName
		c.TplName = "login.html"
		return
	}
	if user.Password != password {
		c.Data["errMsg"] = "用户名或密码错误"
		c.Data["userName"] = userName
		c.TplName = "login.html"
		return
	}
	if user.Active != true {
		c.Data["errMsg"] = "用户名未激活，请先前往邮箱激活"
		c.Data["userName"] = userName
		c.TplName = "login.html"
		return
	}
	// 记住用户名处理
	rememberMe := c.GetString("remember_me")
	if rememberMe == "on" {
		// base64加密
		tempUserName := base64.StdEncoding.EncodeToString([]byte(userName))
		c.Ctx.SetCookie("userName", tempUserName, time.Second * 3600 * 24 * 3)
	} else {
		c.Ctx.SetCookie("userName", userName, -1)
	}
	// 返回视图
	c.Redirect("/", 302)
}
// 显示主页
func (c *UserController) ShowIndex(){
	// 返回视图
	c.TplName = "index.html"
}
