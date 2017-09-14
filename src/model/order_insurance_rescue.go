package model

import (
	proto_order "ocenter/src/proto_order"
	gen "ocenter/src/proto_order/insurance"

	"database/sql"
	"ocenter/src/cst"
	"ocenter/src/db"
	"ocenter/src/model/config"
	"ocenter/src/model/message"
	"time"
)

// InsuranceRescueOrder 保险救援
type InsuranceRescueOrder struct {
	*Order
}

// NewInsuranceRescueOrder 新建保险救援订单
func NewInsuranceRescueOrder() *InsuranceRescueOrder {
	return &InsuranceRescueOrder{Order: NewOrder(cst.OrderTypeInsuranceRescue)}
}

// LoadInsuranceRescueOrder 加载保险救援订单
func LoadInsuranceRescueOrder(id int64) (*InsuranceRescueOrder, error) {
	o := &InsuranceRescueOrder{Order: &Order{OrderSchema: &OrderSchema{}}}
	o.ID = id
	if err := o.load(); err != nil {
		return nil, err
	}
	return o, nil
}

// TriggerPlaceOrderBy 触发下单
func (o *InsuranceRescueOrder) TriggerPlaceOrderBy(d *Driver, r *Repair, fsm *InsuranceRescueFSM, ro *Org, irr *config.InsuranceRescueRule, args *gen.InsuranceRescueArgs) error {
	o.setDriver(d)
	o.setRepair(r)

	fsm.Driver = d
	fsm.Repair = r
	// 触发订单状态机的创建
	if err := fsm.Event(d.PlaceOrder()); err != nil {
		return err
	}

	o.Amount = NewInsuranceResuceCalculator(ro, irr).Calculate()
	o.setState(fsm.Current())
	o.setLocation(args.GetLocation())
	o.DriverTruckID = d.GetCurrentTruck().GetID()
	o.PlateNumber = d.GetCurrentTruck().PlateNumber
	o.Distance = args.Distance
	o.Insurance = args.Insurer
	o.OrderLogs = o.newLog(d.ID)

	if err := o.save(args.Imgs); err != nil {
		return err
	}
	o.sendOrder(d)
	return nil
}

func (o *InsuranceRescueOrder) save(is []string) error {
	return db.TxExec(func(tx *sql.Tx) error {
		// 保存订单
		if err := o.insert(tx); err != nil {
			return err
		}

		oi := NewOrderImage(o.ID)
		oi.SetImg(is)
		if err := oi.Upload(tx); err != nil {
			return err
		}

		// 保存订单日志
		return o.OrderLogs[0].Save(tx, o.ID)
	})
}

func (o *InsuranceRescueOrder) sendOrder(d *Driver) {
	msg := &message.InsuranceRescueMsg{
		OrderID:     o.ID,
		DriverName:  d.Name,
		PlateNumber: o.PlateNumber,
		CreateTime:  o.CreateTime,
		Distance:    o.Distance,
		OrderType:   o.Type,
		OrderStatus: o.State,
	}
	go message.Sender.Send(o.RepairID, msg)
}

// TriggerDriverCancel 订单取消
func (o *InsuranceRescueOrder) TriggerDriverCancel(d *Driver, fsm *InsuranceRescueFSM, args *proto_order.OrderCancelArgs) error {
	if err := o.cancel(d.Cancel(), d.ID, fsm, args.GetProcessBase().GetLocation()); err != nil {
		return err
	}

	msg := &message.OrderCancelMsg{OrderID: o.ID, OrderStatus: o.State}
	go message.Sender.Send(o.RepairID, msg)
	return nil
}

// TriggerRejectBy 订单取消
func (o *InsuranceRescueOrder) TriggerRejectBy(r *Repair, fsm *InsuranceRescueFSM, args *proto_order.OrderCancelArgs) error {
	if err := o.cancel(r.Cancel(), r.ID, fsm, args.GetProcessBase().GetLocation()); err != nil {
		return err
	}

	msg := &message.OrderCancelMsg{OrderID: o.ID, OrderStatus: o.State}
	go message.Sender.Send(o.DriverID, msg)
	return nil
}

func (o *InsuranceRescueOrder) cancel(event string, operator int64, fsm *InsuranceRescueFSM, location *proto_order.Location) error {
	if o.inPay() {
		return ErrOrderInPay
	}

	if err := fsm.Event(event); err != nil {
		return err
	}

	o.modify(fsm.Current(), location, time.Now())

	o.OrderLogs = o.newLog(operator)
	return o.updateBase()
}

// TriggerPayBy 司机付款
func (o *InsuranceRescueOrder) TriggerPayBy(r *Repair, fsm *InsuranceRescueFSM, args *proto_order.OrderPayArgs) error {

	if err := fsm.Event(r.Pay()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())

	o.OrderLogs = o.newLog(r.ID)
	return o.updatePay()
}

func (o *InsuranceRescueOrder) inPay() bool {
	var c int64
	sqlStr := `SELECT COUNT(1) FROM log_pays WHERE order_id = ?`
	_ = db.DbSource.QueryRow(sqlStr, o.ID).Scan(&c)
	return c != 0
}
