# 鲜淘驿站

基于 **B2C** 的鲜淘驿站，Logo已授权。

鲜淘驿站后台管理系统代码请到这里查看：[https://github.com/lpg-it/xian-tao-admin](https://github.com/lpg-it/xian-tao-admin)

采用 **B/S** 结构，即 **Browser/Server** （浏览器/服务器）结构，构建的一个web网站商城系统。

> `B2C` 是 **Business-to-Customer** 的缩写，而其中文简称为“商对客”。“商对客”是电子商务的一种模式，也就是通常说的直接面向消费者销售产品和服务商业零售模式。这种形式的电子商务一般以网络零售业为主，主要借助于互联网开展在线销售活动。`B2C` 即企业通过互联网为消费者提供一个新型的购物环境——网上商店，消费者通过网络在网上购物、网上支付等消费行为。 

### 安装

- 安装FastDFS、nginx：
- 安装Redis
- 

### 技术栈

- 语言：Golang 1.14.2
- 框架：Beego
- 数据库：MySQL、Redis
- 分布式文件存储：FastDFS
- web服务器配置：Nginx

### 主体模块及实现功能

具体请点击这里查看：[鲜淘驿站-主体模块及实现功能](https://github.com/lpg-it/xian-tao/blob/master/doc/鲜淘驿站-主体模块及实现功能.md)

### SKU与SPU概念

​		**SPU = Standard Product Unit** (标准产品单位)

​		SPU 是商品信息聚合的最小单位，是一组可复用、易检索的标准化信息的集合，该集合描述 了一个产品的特性。通俗点讲，属性值、特性相同的商品就可以称为一个 SPU。 

​		例如：iphone8 就是一个 SPU，与商家，与颜色、款式、套餐都无关。

​		**SKU=stock keeping unit**(库存量单位)

​		SKU 即库存进出计量的单位， 可以是以件、盒、托盘等为单位。 

​		SKU 是物理上不可分割的最小存货单元。在使用时要根据不同业态，不同管理模式来处理。 在服装、鞋类商品中使用最多最普遍。 

​		例如：纺织品中一个 SKU 通常表示:规格、颜色、款式。

### 数据库表

具体请点击这里查看：[鲜淘驿站-数据库设计](https://github.com/lpg-it/xian-tao/blob/master/doc/鲜淘驿站-数据库设计.md)



该项目是一个完整的电商项目流程

PS：**此项目纯属个人学习项目**。