package message

import "time"

// Message 消息接口
type Message interface {
	Marshal() []byte
}

// RescueOrderMsg 救援订单消息
type RescueOrderMsg struct {
	OrderID        int64
	Distance       float32
	Subsidy        string
	PlaceOrderType int32
}

// Marshal 序列化
func (r *RescueOrderMsg) Marshal() []byte {
	//TODO
	// protobuf 或者 普通json转换
	return []byte{0x1}
}

// InsuranceRescueMsg 保险救援订单消息
type InsuranceRescueMsg struct {
	OrderID     int64
	Distance    float32
	DriverName  string
	PlateNumber string
	CreateTime  time.Time
	OrderType   int32
	OrderStatus int32
	ExpireTime  int32
}

// Marshal 序列化
func (i *InsuranceRescueMsg) Marshal() []byte {
	//TODO
	// protobuf 或者 普通json转换
	return []byte{0x1}
}

// OrderCancelMsg 订单取消消息
type OrderCancelMsg struct {
	OrderID     int64
	OrderStatus int32
}

// Marshal 序列化
func (o *OrderCancelMsg) Marshal() []byte {
	//TODO
	// protobuf 或者 普通json转换
	return []byte{0x1}
}
