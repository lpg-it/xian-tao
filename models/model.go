package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// 用户表
type User struct {
	Id        int
	Name      string       `orm:"size(20);unique"` // 用户名
	Password  string       `orm:"size(20)"`        // 密码
	Email     string       `orm:"size(50)"`        // 邮箱
	Active    bool         `orm:"default(false)"`  // 是否激活，默认未激活
	Power     int          `orm:"default(0)"`      // 权限设置，0表示普通用户，1表示管理员
	Address   []*Address   `orm:"reverse(many)"`
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

// 地址表
type Address struct {
	Id        int
	Receiver  string       `orm:"size(20)"`       // 收件人
	Addr      string       `orm:"size(50)"`       // 收件地址
	ZipCode   string       `orm:"size(20)"`       // 邮编
	Phone     string       `orm:"size(20)"`       // 联系方式
	IsDefault bool         `orm:"default(false)"` // 是否是默认地址
	User      *User        `orm:"rel(fk)"`        // 用户ID
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

// 商品SPU表
type Goods struct {
	Id       int
	Name     string      `orm:"size(50)"`  // 商品名称
	Detail   string      `orm:"size(200)"` // 商品详细描述
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
}

// 商品类型表
type GoodsType struct {
	Id                   int
	Name                 string                  `orm:"size(20)"` // 类型名称
	Logo                 string                  // 类型Logo
	Image                string                  // 类型图片
	GoodsSKU             []*GoodsSKU             `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"` // 首页分类商品
}

// 商品图片表
type GoodsImage struct {
	Id       int
	Image    string    // 商品图片
	GoodsSKU *GoodsSKU `orm:"rel(fk)"` // 商品SKU
}

// 商品SKU表
type GoodsSKU struct {
	Id                   int
	Goods                *Goods                  `orm:"rel(fk)"`   // 商品SPU
	GoodsType            *GoodsType              `orm:"rel(fk)"`   // 商品类型
	Name                 string                  `orm:"size(50)"`  // 商品名称
	Desc                 string                  `orm:"size(100)"` // 商品简介
	Price                int                     // 商品价格
	Unite                string                  `orm:"size(20)"` // 商品单位
	Image                string                  // 商品图片
	Stock                int                     `orm:"default(1)"`                  // 商品库存
	Sales                int                     `orm:"default(0)"`                  // 商品销量
	Status               int                     `orm:"default(1)"`                  // 商品状态：是否有效，默认有效
	Time                 time.Time               `orm:"auto_now_add"` // 添加时间
	GoodsImage           []*GoodsImage           `orm:"reverse(many)"`               // 商品图片
	IndexGoodsBanner     []*IndexGoodsBanner     `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"`
	OrderGoods           []*OrderGoods           `orm:"reverse(many)"`
}

// 首页轮播商品展示表
type IndexGoodsBanner struct {
	Id       int
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
	Image    string    // 商品图片
	Index    int       `orm:"default(0)"` // 商品展示顺序

}

// 首页分类商品展示表
type IndexTypeGoodsBanner struct {
	Id          int
	GoodsType   *GoodsType `orm:"rel(fk)"`
	GoodsSKU    *GoodsSKU  `orm:"rel(fk)"`
	DisplayType int        `orm:"default(1)"` // 展示类型：0 代表文字；1 代表图片
	Index       int        `orm:"default(0)"` // 展示顺序
}

// 首页促销商品展示表
type IndexPromotionBanner struct {
	Id    int
	Name  string `orm:"size(20)"` // 活动名称
	Url   string `orm:"size(50)"` // 活动链接
	Image string // 活动图片
	Index int    `orm:"default(0)"` // 展示顺序
}

// 订单商品表
type OrderGoods struct {
	Id        int
	OrderInfo *OrderInfo `orm:"rel(fk)"`
	GoodsSKU  *GoodsSKU  `orm:"rel(fk)"`
	Count     int        `orm:"default(1)"` // 商品数量
	Price     int        // 商品价格
	Comment   string     `orm:"default('');size(200)"` // 评论内容
}

// 订单表
type OrderInfo struct {
	Id           int
	OrderId      string        `orm:"size(20);unique"` // 订单号
	User         *User         `orm:"rel(fk)"`
	Address      *Address      `orm:"rel(fk)"` // 收货地址
	PayMethod    int           // 支付方式
	TotalCount   int           `orm:"default(1)"` // 商品数量
	TotalPrice   int           // 商品总价（包含运费）
	TransitPrice int           `orm:"default(0)"`                  // 运费
	OrderStatus  int           `orm:"default(1)"`                  // 支付状态: 1未支付， 2支付成功
	TradeNo      string        `orm:"size(20);default('')"`        // 支付编号
	Time         time.Time     `orm:"auto_now_add"` // 时间
	OrderGoods   []*OrderGoods `orm:"reverse(many)"`
}

func init() {
	// 注册默认数据库
	orm.RegisterDataBase("default", "mysql", "root:数据库密码%%@tcp(127.0.0.1:3306)/xian_tao?charset=utf8")

	// 注册 model
	orm.RegisterModel(new(User), new(Address), new(Goods), new(GoodsType), new(GoodsImage), new(GoodsSKU), new(IndexGoodsBanner), new(IndexTypeGoodsBanner), new(IndexPromotionBanner), new(OrderGoods), new(OrderInfo))

	// 创建表
	orm.RunSyncdb("default", false, true)
}
