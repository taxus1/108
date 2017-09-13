package model

import (
	"database/sql"
	"ocenter/src/cst"
	"ocenter/src/db"
	"time"

	proto_order "ocenter/src/proto_order"
	gen "ocenter/src/proto_order/repair"
)

//RepairOrder 订单对象
type RepairOrder struct {
	*Order
}

// NewRepairOrder 创建订单对象
func NewRepairOrder() *RepairOrder {
	return &RepairOrder{Order: NewOrder(cst.OrderTypeTruckRepair)}
}

// LoadRepairOrder 加载已有订单
func LoadRepairOrder(id int64) (*RepairOrder, error) {
	o := &RepairOrder{}
	o.ID = id
	err := o.load()
	return o, err
}

// TriggerPlaceOrderBy 触发下单
func (o *RepairOrder) TriggerPlaceOrderBy(d *Driver, r *Repair, fsm *RepairFSM, args *gen.RepairCreateArgs) error {
	fsm.Driver = d
	fsm.Repair = r
	if err := fsm.Event(d.PlaceOrder()); err != nil {
		return err
	}
	o.setState(fsm.Current())

	o.OrderInfos = NewOrderInfos(args.GetOrderInfos(), o.Type)
	if err := o.setAmount(NewRepairCalculator(o.OrderInfos, args.Extra, 0)); err != nil {
		return err
	}

	o.SubType = o.OrderInfos.subTypes()
	o.Distance = o.OrderInfos.FindMile()
	o.DriverTruckID = d.GetCurrentTruck().GetID()
	o.setDriver(d)
	o.setRepair(r)
	o.setLocation(args.GetLocation())
	o.setCustomerPrice(args.Extra, 0)
	o.OrderLogs = o.newLog(d.ID)

	return o.saveRelative(r.SetToBusy)
}

// TriggerBargainBy 触发订单议价
func (o *RepairOrder) TriggerBargainBy(r *Repair, fsm *RepairFSM, args *gen.OrderBargainArgs) error {
	if err := fsm.Event(r.Bargain()); err != nil {
		return err
	}

	o.OrderInfos = NewOrderInfos(args.OrderInfos, o.Type)
	if err := o.setAmount(NewRepairCalculator(o.OrderInfos, args.CustomerPrice.Extra, args.CustomerPrice.AccessoriesFee)); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(r.ID)
	return o.updateBargain()
}

// 更新订单，插入订单材料，创建订单日志
func (o *RepairOrder) updateBargain() error {
	return db.TxExec(func(tx *sql.Tx) error {
		// 保存订单
		return o.updatePriceChange(tx)
	})
}

// TriggerConfirmBy 触发订单确认
func (o *RepairOrder) TriggerConfirmBy(d *Driver, fsm *RepairFSM, args *gen.OrderConfirmeArgs) error {
	if err := fsm.Event(d.Confirm()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(d.ID)
	return o.updateConfirm()
}

func (o *RepairOrder) updateConfirm() error {
	return o.updateBase()
}

// TriggerRunBy 触发技工出发中
func (o *RepairOrder) TriggerRunBy(r *Repair, fsm *RepairFSM, args *gen.OrderRunArgs) error {
	if err := fsm.Event(r.Run()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(r.ID)
	return o.updateRun()
}

func (o *RepairOrder) updateRun() error {
	return o.updateBase()
}

// TriggerHandBy 触发技工到达修理中
func (o *RepairOrder) TriggerHandBy(r *Repair, fsm *RepairFSM, args *gen.OrderHandArgs) error {
	if err := fsm.Event(r.Hand()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(r.ID)
	return o.updateHand()
}

func (o *RepairOrder) updateHand() error {
	return o.updateBase()
}

// TriggerFinishBy 触发技工完成修理
func (o *RepairOrder) TriggerFinishBy(r *Repair, fsm *RepairFSM, args *gen.OrderFinishArgs) error {
	if err := fsm.Event(r.Finish()); err != nil {
		return err
	}

	o.OrderInfos = o.newInfo(args.GetOrderInfos())
	if err := o.setAmount(NewRepairCalculator(o.OrderInfos, args.GetCustomerPrice().Extra, args.GetCustomerPrice().AccessoriesFee)); err != nil {
		return err
	}
	o.setCustomerPrice(args.GetCustomerPrice().Extra, args.GetCustomerPrice().AccessoriesFee)

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(r.ID)
	return o.updateFinish(r)
}

func (o *RepairOrder) updateFinish(r *Repair) error {
	return db.TxExec(func(tx *sql.Tx) error {
		if err := o.updatePriceChange(tx); err != nil {
			return err
		}

		// 技工设为接单中
		return r.SetToWork(tx)
	})
}

// TriggerDriverCancel 订单取消
func (o *RepairOrder) TriggerDriverCancel(d *Driver, r *Repair, fsm *RepairFSM, args *proto_order.OrderCancelArgs) error {
	return o.cancel(d.Cancel(), d.ID, r, fsm, args)
}

// TriggerRepairCancel 订单取消
func (o *RepairOrder) TriggerRepairCancel(r *Repair, fsm *RepairFSM, args *proto_order.OrderCancelArgs) error {
	return o.cancel(r.Cancel(), r.ID, r, fsm, args)
}

func (o *RepairOrder) cancel(event string, operator int64, r *Repair, fsm *RepairFSM, args *proto_order.OrderCancelArgs) error {
	if err := fsm.Event(event); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(operator)
	return o.updateCancel(r)
}

func (o *RepairOrder) updateCancel(r *Repair) error {
	return db.TxExec(func(tx *sql.Tx) error {
		if err := o.updateState(tx); err != nil {
			return err
		}

		// 技工设为接单中
		return r.SetToWork(tx)
	})
}

func (o *RepairOrder) updatePriceChange(tx *sql.Tx) error {
	if err := o.updatePrice(tx); err != nil {
		return err
	}

	// 保存订单日志
	return o.saveInfoAndLog(tx)
}

func (o *RepairOrder) setCustomerPrice(e, af int64) {
	o.Extra = e
	o.AccessoriesFee = af
}
