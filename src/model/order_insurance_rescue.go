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
	InsuranceRescueFSM *InsuranceRescueFSM
}

// NewInsuranceRescueOrder 新建保险救援订单
func NewInsuranceRescueOrder() *InsuranceRescueOrder {
	return &InsuranceRescueOrder{Order: NewOrder(cst.OrderTypeInsuranceRescue)}
}

// TriggerPlaceOrderBy 触发下单
func (o *InsuranceRescueOrder) TriggerPlaceOrderBy(d *Driver, r *Repair, args *gen.InsuranceRescueArgs) error {
	o.setDriver(d)
	o.setRepair(r)

	o.InsuranceRescueFSM = NewInsuranceRescueFSM(o)

	// 触发订单状态机的创建
	if err := o.InsuranceRescueFSM.Event(d.PlaceOrder()); err != nil {
		return err
	}

	ro, err := LoadOrg(o.RepairOrgID)
	if err != nil {
		return err
	}

	irr, err := config.LoadInsuranceRescueRule()
	if err != nil {
		return err
	}
	c := NewInsuranceResuceCalculator(ro, irr)
	o.Amount = c.Calculate()

	o.setState(o.InsuranceRescueFSM.Current())
	o.setLocation(args.GetLocation())

	o.DriverTruckID = d.GetCurrentTruck().GetID()
	o.PlateNumber = d.GetCurrentTruck().PlateNumber
	o.SubType = o.OrderInfos.subTypes()
	o.Distance = o.OrderInfos.FindMile()
	o.Insurance = args.Insurer
	o.newLog(d.ID)

	if err := o.save(args.Imgs); err != nil {
		return err
	}

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
		return o.saveInfoAndLog(tx)
	})
}

// TriggerDriverCancel 订单取消
func (o *InsuranceRescueOrder) TriggerDriverCancel(d *Driver, args *proto_order.OrderCancelArgs) error {
	if err := o.cancel(d.Cancel(), d.ID, args.GetProcessBase().GetLocation()); err != nil {
		return err
	}

	msg := &message.OrderCancelMsg{OrderID: o.ID, OrderStatus: o.State}
	go message.Sender.Send(o.RepairID, msg)
	return nil
}

// TriggerRejectBy 订单取消
func (o *InsuranceRescueOrder) TriggerRejectBy(r *Repair, args *proto_order.OrderCancelArgs) error {
	if err := o.cancel(r.Cancel(), r.ID, args.GetProcessBase().GetLocation()); err != nil {
		return err
	}

	msg := &message.OrderCancelMsg{OrderID: o.ID, OrderStatus: o.State}
	go message.Sender.Send(o.DriverID, msg)
	return nil
}

func (o *InsuranceRescueOrder) cancel(event string, operator int64, location *proto_order.Location) error {
	if o.inPay() {
		return ErrOrderInPay
	}

	if err := o.InsuranceRescueFSM.Event(event); err != nil {
		return err
	}

	o.modify(
		o.InsuranceRescueFSM.Current(),
		location,
		time.Now(),
	)

	o.newLog(operator)
	return o.updateBase()
}

// TriggerPayBy 司机付款
func (o *InsuranceRescueOrder) TriggerPayBy(r *Repair, args *proto_order.OrderPayArgs) error {

	if err := o.InsuranceRescueFSM.Event(r.Pay()); err != nil {
		return err
	}

	o.modify(
		o.InsuranceRescueFSM.Current(),
		args.GetProcessBase().GetLocation(),
		time.Now(),
	)

	o.newLog(r.ID)
	return o.updatePay()
}

// TriggerAssessBy 司机评价
func (o *InsuranceRescueOrder) TriggerAssessBy(d *Driver, args *proto_order.OrderAssessArgs) error {
	if err := o.InsuranceRescueFSM.Event(d.Assess()); err != nil {
		return err
	}

	o.modify(
		o.InsuranceRescueFSM.Current(),
		args.GetProcessBase().GetLocation(),
		time.Now(),
	)

	o.newLog(d.ID)
	return o.updateAccess()
}

func (o *InsuranceRescueOrder) inPay() bool {
	var c int64
	sqlStr := `SELECT COUNT(1) FROM log_pays WHERE order_id = ?`
	_ = db.DbSource.QueryRow(sqlStr, o.ID).Scan(&c)
	return c != 0
}
