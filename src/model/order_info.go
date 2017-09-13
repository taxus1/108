package model

import (
	"database/sql"
	"ocenter/src/db"
	proto_order "ocenter/src/proto_order"
	"ocenter/src/util"

	"strconv"
	"strings"
)

// orderInfoModel 订单关联的物品
type orderInfoModel Model

var colums = []string{
	"order_id",
	"goods_id",
	"goods_name",
	"goods_code",
	"price",
	"amount",
	"parent",
	"order_process",
	"order_type",
	"info_type",
}

// OrderInfoSchema 订单信息表
type OrderInfoSchema struct {
	ID int64

	//物品名字
	Name string

	//数量
	Amount float32

	//价格
	Price int64

	//订单ID
	OrderID int64

	//物品ID
	GoodsID int64

	//维修订单物品上级
	ParentID int64

	//订单进度
	OrderProcess int64

	//部件类型
	UnitType int32

	//订单信息类型
	InfoType int32

	//物品编码
	GoodsCode string
}

// OrderInfo 订单关联的物品
type OrderInfo struct {
	*OrderInfoSchema
}

// OrderInfos 一般都会出现多个，所以定义数组类型的
type OrderInfos []*OrderInfo

// NewOrderInfos 创建方法
func NewOrderInfos(args []*proto_order.OrderInfo, orderType int32) *OrderInfos {
	infos := make(OrderInfos, len(args))
	for i, v := range args {
		info := &OrderInfo{
			&OrderInfoSchema{
				Name:      v.Name,
				GoodsID:   v.GoodsID,
				Price:     v.Price,
				Amount:    v.Amount,
				ParentID:  v.ParentID,
				InfoType:  v.InfoType,
				GoodsCode: v.GoodsCode,
			},
		}
		infos[i] = info
		info.setUnitType(v.UnitType, orderType)
	}
	return &infos
}

// FindMile 从订单物品中获取里程数量
func (ois *OrderInfos) FindMile() float32 {
	for _, v := range *ois {
		if v.GoodsCode == "TIRE_MILE" {
			return v.Amount
		}
	}
	return 0
}

// SubTypes 从订单物品中子类型
func (ois *OrderInfos) subTypes() *orderInfoTypes {
	var types orderInfoTypes
	for _, v := range *ois {
		types.appendType(v.getSubType())
	}
	return &types
}

// MaterialMoney 从订单物品中订单材料总金额
func (ois *OrderInfos) MaterialMoney() (money *Money) {
	m := Money(0)
	for _, v := range *ois {
		m.add(v.getItemMoney())
	}
	return &m
}

// BatchInsert 批量插入
func (ois *OrderInfos) BatchInsert(tx *sql.Tx, orderID int64, state int32) error {
	sqlStr := db.BuildInsert("order_info", colums, len(*ois))
	stmt, _ := tx.Prepare(sqlStr)

	values := ois.eachValues(orderID, state)
	_, err := stmt.Exec(values...)
	return err
}

func (ois *OrderInfos) eachValues(oderID int64, state int32) []interface{} {
	values := []interface{}{}
	for _, info := range *ois {
		values = append(values, oderID)
		values = append(values, info.GoodsID)
		values = append(values, info.Name)
		values = append(values, info.GoodsCode)
		values = append(values, info.Price)
		values = append(values, info.Amount)
		values = append(values, info.ParentID)
		values = append(values, state)
		values = append(values, info.UnitType)
		values = append(values, info.InfoType)
	}
	return values
}

func (oi *OrderInfo) getSubType() string {
	return strconv.Itoa(int(oi.UnitType))
}

//获取订单单个条目的总价 数量 * 单价
func (oi *OrderInfo) getItemMoney() float32 {
	if oi.Amount == 0 || oi.Price == 0 {
		return 0
	}
	return oi.Amount * (float32(oi.Price) / 100)
}

func (oi *OrderInfo) setUnitType(unitType, orderType int32) {
	if unitType == 0 {
		oi.UnitType = orderType
	} else {
		oi.UnitType = unitType
	}
}

// 订单类型集合
type orderInfoTypes []string

func (oit *orderInfoTypes) join(sep string) string {
	return strings.Join(*oit, sep)
}

func (oit *orderInfoTypes) appendType(unitType string) {
	if !util.SearchString(*oit, unitType) {
		*oit = append(*oit, unitType)
	}
}
