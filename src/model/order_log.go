package model

import (
	"database/sql"
	"ocenter/src/db"

	"time"
)

// orderLogModel 订单日志
type orderLogModel Model

var fileds = []string{
	"order_id",
	"status",
	"user_id",
	"driver_id",
	"repair_id",
	"amount",
	"lat",
	"lng",
	"order_address",
	"extra",
	"distance",
	"miles",
	"accessories_fee",
	"oil_gun",
	"create_time",
}

// OrderLogSchema 订单日志表
type OrderLogSchema struct {
	ID int64

	//订单ID
	OrderID int64

	//状态标识
	OrderState int32

	//司机
	DriverID int64

	//技工
	RepairID int64

	//产生这个日志的用户ID
	Operator int64

	//金额
	Amount int64

	//附加费
	Extra int64

	//距离
	Distance float32

	//里程
	Miles int64

	//辅料费
	AccessoriesFee int64

	//坐标lat
	Lat float32

	//坐标lat
	Lng float32

	//订单中文位置
	OrderAddress string

	//油枪编号
	OilGun string

	//处理时间
	CreateTime time.Time
}

// OrderLog 订单日志
type OrderLog struct {
	*OrderLogSchema
}

// NewOrderLog 订单日志对象
func NewOrderLog(order *Order, operator int64) *OrderLog {
	return &OrderLog{
		&OrderLogSchema{
			OrderID:        order.ID,
			Operator:       operator,
			OrderState:     order.State,
			DriverID:       order.DriverID,
			RepairID:       order.RepairID,
			Amount:         order.Amount,
			Extra:          order.Extra,
			Distance:       order.Distance,
			AccessoriesFee: order.AccessoriesFee,
			Lat:            order.Lat,
			Lng:            order.Lng,
			OrderAddress:   order.Address,
			OilGun:         order.OilGun,
			CreateTime:     time.Now(),
		},
	}
}

// Save 创建订单日志
func (l *OrderLog) Save(tx *sql.Tx, orderID int64) error {
	sqlStr := db.BuildInsert("order_log", fileds, 1)

	values := []interface{}{
		orderID,
		l.OrderState,
		l.Operator,
		l.DriverID,
		l.RepairID,
		l.Amount,
		l.Lat,
		l.Lng,
		l.OrderAddress,
		l.Extra,
		l.Distance,
		l.Miles,
		l.AccessoriesFee,
		l.OilGun,
		l.CreateTime,
	}
	_, err := db.SaveTx(tx, sqlStr, values)

	return err
}
