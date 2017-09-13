package model

import (
	proto_order "ocenter/src/proto_order"
	gen "ocenter/src/proto_order/oil"

	"ocenter/src/cst"
	"time"
)

//OilOrder 订单对象
type OilOrder struct {
	*Order
}

// NewOilOrder 创建订单对象
func NewOilOrder() *OilOrder {
	return &OilOrder{Order: NewOrder(cst.OrderTypeOilStation)}
}

// LoadOilOrder 加载已有订单
func LoadOilOrder(id int64) (*OilOrder, error) {
	o := &OilOrder{&Order{OrderSchema: &OrderSchema{}}}
	o.ID = id
	err := o.load()
	return o, err
}

// TriggerPlaceOrderByDriver 触发下单
func (o *OilOrder) TriggerPlaceOrderByDriver(d *Driver, r *Repair, fsm *OilFSM, args *gen.OilCreateArgs) error {
	fsm.Driver = d
	if err := fsm.Event(d.PlaceOrder()); err != nil {
		return err
	}
	o.setState(fsm.Current())

	o.OrderInfos = o.newInfo([]*proto_order.OrderInfo{args.GetOrderInfo()})
	if err := o.setAmount(NewOilCalculator(o.OrderInfos)); err != nil {
		return err
	}

	o.SubType = o.OrderInfos.subTypes()
	o.DriverTruckID = d.GetCurrentTruck().GetID()
	o.setDriver(d)
	o.setRepair(r)
	o.setLocation(args.GetLocation())
	o.OrderLogs = o.newLog(d.ID)

	return o.saveRelative()
}

// TriggerPlaceOrderByOiler 触发下单
func (o *OilOrder) TriggerPlaceOrderByOiler(r *Repair, org *Org, fsm *OilFSM, args *gen.OilCreateArgs) error {
	if err := fsm.Event(r.PlaceOrder()); err != nil {
		return err
	}
	o.setState(fsm.Current())

	o.OrderInfos = o.newInfo([]*proto_order.OrderInfo{args.GetOrderInfo()})
	if err := o.setAmount(NewOilCalculator(o.OrderInfos)); err != nil {
		return err
	}

	o.ServiceFee = org.ServiceFee(args.GetOrderInfo().Amount, o.Amount)
	o.SubType = o.OrderInfos.subTypes()
	o.setLocation(args.GetLocation())
	o.OrderLogs = o.newLog(r.ID)

	return o.saveRelative()
}

// TriggerDriverClose 订单关闭
func (o *OilOrder) TriggerDriverClose(d *Driver, fsm *OilFSM, args *proto_order.OrderCancelArgs) error {
	if err := fsm.Event(d.Close()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(d.ID)
	return o.updateClose()
}

func (o *OilOrder) updateClose() error {
	return o.updateBase()
}
