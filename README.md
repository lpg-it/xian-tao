# 鲜淘驿站

基于 **B2C** 的鲜淘驿站，Logo已授权。

鲜淘驿站后台管理系统代码请到这里查看：[https://github.com/lpg-it/xian-tao-admin](https://github.com/lpg-it/xian-tao-admin)

采用 **B/S** 结构，即 **Browser/Server** （浏览器/服务器）结构，构建的一个web网站商城系统。

> `B2C` 是 **Business-to-Customer** 的缩写，而其中文简称为“商对客”。“商对客”是电子商务的一种模式，也就是通常说的直接面向消费者销售产品和服务商业零售模式。这种形式的电子商务一般以网络零售业为主，主要借助于互联网开展在线销售活动。`B2C` 即企业通过互联网为消费者提供一个新型的购物环境——网上商店，消费者通过网络在网上购物、网上支付等消费行为。 

### 技术栈

- 语言：Golang1.14.2
- 框架：Beego
- 数据库：MySQL、Redis
- 分布式文件存储：FastDFS
- web服务器配置：Nginx

### 主体模块及实现功能

- 用户模块
  - **注册页**
    - 注册时校验用户是否已被注册。
    - 给注册用户发送激活邮件，用户点击邮件中的链接完成用户激活。
    - 完成用户信息的注册。
  - **登录页**
    - 实现用户的登录功能。
  - **用户中心**
    - 用户中心信息页：显示登录用户的信息，包括用户名、电话和地址，**同时页面下方显示出用户最近浏览的商品信息（Redis存储）**。
    - 用户中心地址页：显示登录用户的默认收件地址，页面下方的表单可以新增用户的收货地址。
    - 用户中心订单页：显示登录用户的订单信息。
  - **其它**
    - 如果用户已登录，页面顶部显示登录用户的信息。
- 商品模块
  - **首页**
    - 动态指定首页轮播商品信息
    - 动态指定首页活动信息
    - 动态获取商品的种类信息并显示
    - 动态指定首页显示的每个种类的商品（包括文字商品和图片商品）
    - 点击某一个商品时跳转到商品的详情页面
  - **商品详情页**
    - 显示出某个商品的详细信息
    - 页面的左下方显示出该种类的2个新上架商品信息
- 购物车模块
- 订单模块
- 后台模块



### SKU与SPU概念

​		**SPU = Standard Product Unit** (标准产品单位)

​		SPU 是商品信息聚合的最小单位，是一组可复用、易检索的标准化信息的集合，该集合描述 了一个产品的特性。通俗点讲，属性值、特性相同的商品就可以称为一个 SPU。 

​		例如：iphone8 就是一个 SPU，与商家，与颜色、款式、套餐都无关。

​		**SKU=stock keeping unit**(库存量单位)

​		SKU 即库存进出计量的单位， 可以是以件、盒、托盘等为单位。 

​		SKU 是物理上不可分割的最小存货单元。在使用时要根据不同业态，不同管理模式来处理。 在服装、鞋类商品中使用最多最普遍。 

​		例如：纺织品中一个 SKU 通常表示:规格、颜色、款式。

### 数据库表

#### MySQL

##### 用户表：User

| 字段          | 备注                             |
| ------------- | -------------------------------- |
| Id            |                                  |
| Name          | 用户名                           |
| Password      | 密码                             |
| Email         | 邮箱                             |
| Active        | 是否激活                         |
| Power         | 权限设置：0表示未激活，1表示激活 |
| **Address**   |                                  |
| **OrderInfo** |                                  |

##### 地址表：Address

| 字段          | 备注         |
| ------------- | ------------ |
| Id            |              |
| Receiver      | 收件人       |
| Addr          | 收件地址     |
| ZipCode       | 邮编         |
| Phone         | 联系方式     |
| IsDefault     | 是否默认地址 |
| **User**      | 用户ID       |
| **OrderInfo** |              |

##### 商品SPU表：Goods

| 字段         | 备注         |
| ------------ | ------------ |
| Id           |              |
| Name         | 商品名称     |
| Detail       | 商品详细描述 |
| **GoodsSKU** |              |

##### 商品类型表：GoodsType

| 字段                     | 备注     |
| ------------------------ | -------- |
| Id                       |          |
| Name                     | 类型名称 |
| Logo                     | 类型图标 |
| Image                    | 类型图片 |
| **GoodsSKU**             |          |
| **IndexTypeGoodsBanner** |          |

##### 商品图片表：GoodsImage

| 字段         | 备注     |
| ------------ | -------- |
| Id           |          |
| Image        | 商品图片 |
| **GoodsSKU** | 商品SKU  |

##### 商品SKU表：GoodsSKU

| 字段                     | 备注                 |
| ------------------------ | -------------------- |
| ID                       |                      |
| Name                     | 商品名称             |
| Desc                     | 商品简介             |
| Price                    | 商品价格             |
| Unite                    | 商品单位             |
| Image                    | 商品图片             |
| Stock                    | 商品库存             |
| Sales                    | 商品销量             |
| Status                   | 商品状态（是否有效） |
| Time                     | 商品添加时间         |
| **Goods**                | 商品SPU              |
| **GoodsType**            | 商品所属种类         |
| **GoodsImage**           |                      |
| **IndexGoodsBanner**     |                      |
| **IndexTypeGoodsBanner** |                      |
| **OrderGoods**           |                      |

##### 首页轮播商品展示表：IndexGoodsBanner

| 字段         | 备注     |
| ------------ | -------- |
| Id           |          |
| Image        | 商品图片 |
| Index        | 展示顺序 |
| **GoodsSKU** | 商品SKU  |

##### 首页分类商品展示表：IndexTypeGoodsBanner

| 字段          | 备注                           |
| ------------- | ------------------------------ |
| Id            |                                |
| **GoodsType** | 商品类型                       |
| **GoodsSKU**  | 商品SKU                        |
| DisplayType   | 展示类型：0代表图片，1代表文字 |
| Index         | 展示顺序                       |

##### 首页促销商品展示表：IndexPromotionBanner

| 字段  | 备注     |
| ----- | -------- |
| Id    |          |
| Name  | 活动名称 |
| Url   | 活动链接 |
| Image | 活动图片 |
| Index | 展示顺序 |

##### 订单商品表：OrderGoods

| 字段          | 备注     |
| ------------- | -------- |
| Id            |          |
| **OrderInfo** | 订单     |
| **GoodsSKU**  | 商品SKU  |
| Count         | 商品数量 |
| Price         | 商品价格 |
| Comment       | 评论内容 |

##### 订单表：OrderInfo

| 字段           | 备注                     |
| -------------- | ------------------------ |
| Id             |                          |
| OrderId        | 订单号                   |
| **User**       | 用户                     |
| **Address**    | 收货地址                 |
| PayMethod      | 支付方式                 |
| TotalCount     | 商品数量                 |
| TotalPrice     | 商品总价：**空间换时间** |
| TransitPrice   | 运费                     |
| OrderStatus    | 订单状态：已付款/未付款  |
| TradeNo        | 支付编号                 |
| Time           | 评论时间                 |
| **OrderGoods** |                          |

##### MySQL初始化

```go
func init(){
    // 设置默认数据库
    orm.RegisterDataBase("default", "mysql", "root:密码@tcp(127.0.0.1:3306)/xian_tao?charset=utf8")
    
    // 注册 model
    orm.RegisterModel(new(User), new(Address), new(Goods), new(GoodsType), new(GoodsImage), new(GoodsSKU), new(IndexGoodsBanner), new(IndexTypeGoodsBanner), new(IndexPromotionBanner), new(OrderGoods), new(OrderInfo))
    
    // 创建表
    orm.RunSyncdb("default", false, true)
}
```

#### Redis

##### 购物车数据

##### 历史浏览记录



该项目是一个完整的电商项目流程

PS：**此项目纯属个人学习项目**。