package model

import (
	"fmt"
	"ocenter/src/cst"
	"ocenter/src/db"
	"ocenter/src/model/config"
	proto_order "ocenter/src/proto_order"
	gen "ocenter/src/proto_order/rescue"

	"database/sql"
	"time"
)

const (
	placeTypeRepair = 1
	placeTypeRescue = 2
	placeTypeDirect = 3
)

//RepairRescueOrder 订单对象
type RepairRescueOrder struct {
	*RepairOrder

	OrderImage *OrderImage
}

// NewRepairRescueOrder 创建订单对象
func NewRepairRescueOrder() *RepairRescueOrder {
	return &RepairRescueOrder{RepairOrder: NewRepairOrder()}
}

// LoadRepairRescueOrder 加载已有订单
func LoadRepairRescueOrder(id int64) (*RepairRescueOrder, error) {
	o := &RepairRescueOrder{RepairOrder: &RepairOrder{Order: &Order{OrderSchema: &OrderSchema{}}}}
	o.ID = id
	if err := o.load(); err != nil {
		return nil, err
	}
	oi, err := LoadOrderImage(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	o.OrderImage = oi

	return o, nil
}

// TriggerRescueBy 触发下单
func (o *RepairRescueOrder) TriggerRescueBy(d *Driver, fsm *RepairRescueFSM, c *config.OtherConf, args *gen.RepairRescueArgs) error {
	fsm.Driver = d
	o.PlaceOrderType = placeTypeRescue
	if err := o.rescue(d, fsm, c, args); err != nil {
		return err
	}

	if err := o.saveRescue(); err != nil {
		return err
	}

	rs, err := FindRepairor(d, args.GetLocation(), o.SubType)
	if err != nil || len(rs) == 0 {
		return nil
	}

	o.dispatcheTo(rs...)
	return nil
}

func (o *RepairRescueOrder) saveRescue() error {
	return db.TxExec(func(tx *sql.Tx) error {
		// 保存订单
		if err := o.insert(tx); err != nil {
			return err
		}

		rescue := NewRescue(o)
		if err := rescue.insert(tx); err != nil {
			return err
		}

		// 保存订单日志
		return o.saveInfoAndLog(tx)
	})
}

// TriggerDirectRescueBy P2P下单
func (o *RepairRescueOrder) TriggerDirectRescueBy(d *Driver, r *Repair, fsm *RepairRescueFSM, c *config.OtherConf, args *gen.DirectRescueArgs) error {
	o.setRepair(r)
	fsm.Repair = r
	o.PlaceOrderType = placeTypeDirect
	if err := o.rescue(d, fsm, c, args.GetRescueArgs()); err != nil {
		return err
	}

	if err := o.saveRelative(); err != nil {
		return err
	}

	rr := &RescueRepair{ID: r.ID, Distance: args.Distance}
	o.dispatcheTo(rr)
	return nil
}

func (o *RepairRescueOrder) rescue(d *Driver, fsm *RepairRescueFSM, c *config.OtherConf, args *gen.RepairRescueArgs) error {
	o.setDriver(d)

	// 触发订单状态机的创建
	fsm.Driver = d
	fsm.Config = c
	if err := fsm.Event(d.PlaceOrder()); err != nil {
		return err
	}

	o.setState(fsm.Current())
	o.setLocation(args.GetLocation())
	o.DriverTruckID = d.GetCurrentTruck().GetID()
	o.OrderInfos = o.newInfo(args.GetOrderInfos())
	o.SubType = o.OrderInfos.subTypes()
	o.Distance = o.OrderInfos.FindMile()
	o.PlateNumber = d.GetCurrentTruck().PlateNumber
	o.OrderLogs = o.newLog(d.ID)

	return nil
}

// TriggerGrabBy 触发技工抢单
func (o *RepairRescueOrder) TriggerGrabBy(r *Repair, fsm *RepairRescueFSM, c *config.GrabConf, g *Rescue, args *gen.RescueGrabArgs) error {
	o.setRepair(r)
	fsm.Repair = r
	if err := fsm.Event(r.Grab()); err != nil {
		return err
	}

	if err := g.CheckGrabed(c); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetLocation(), time.Now())
	o.Distance = args.Distance
	o.OrderLogs = o.newLog(r.ID)
	return o.updateGrab(g, r)
}

func (o *RepairRescueOrder) updateGrab(g *Rescue, r *Repair) error {
	return db.TxExec(func(tx *sql.Tx) error {
		if err := g.Grab(tx, r.ID); err != nil {
			return err
		}

		if err := o.updateGrabeInfo(tx); err != nil {
			return err
		}

		if err := o.OrderLogs[0].Save(tx, o.ID); err != nil {
			return err
		}

		return r.SetToBusy(tx)
	})
}

// TriggerCompleteBy 技工触发完成修理
func (o *RepairRescueOrder) TriggerCompleteBy(r *Repair, fsm *RepairRescueFSM, args *gen.RescueCompleteArgs) error {
	if err := fsm.Event(r.Finish()); err != nil {
		return err
	}

	o.OrderInfos = o.newInfo(args.GetOrderInfos())
	if err := o.setAmount(NewRepairCalculator(o.OrderInfos, args.GetCustomerPrice().Extra, args.GetCustomerPrice().AccessoriesFee)); err != nil {
		return err
	}
	o.setCustomerPrice(args.GetCustomerPrice().Extra, args.GetCustomerPrice().AccessoriesFee)

	o.modify(fsm.Current(), args.GetLocation(), time.Now())
	o.OrderLogs = o.newLog(r.ID)

	oi := &OrderImage{}
	oi.SetImg(args.Image)
	o.OrderImage = oi
	return o.updateFinish(r)
}

func (o *RepairRescueOrder) updateFinish(r *Repair) error {
	return db.TxExec(func(tx *sql.Tx) error {
		if err := o.updatePriceChange(tx); err != nil {
			return err
		}

		if err := o.OrderImage.update(tx); err != nil {
			return err
		}

		// 技工设为接单中
		return r.SetToWork(tx)
	})
}

// TriggerAcceptBy 触发技工接单
func (o *RepairRescueOrder) TriggerAcceptBy(r *Repair, fsm *RepairRescueFSM, args *gen.RescueAcceptArgs) error {
	o.setRepair(r)
	fsm.Repair = r
	if err := fsm.Event(r.Accept()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetAcceptArgs().GetLocation(), time.Now())
	o.Distance = args.GetAcceptArgs().Distance
	o.OrderLogs = o.newLog(r.ID)
	return o.updateAccept(r)
}

func (o *RepairRescueOrder) updateAccept(r *Repair) error {
	return db.TxExec(func(tx *sql.Tx) error {

		if err := o.updateGrabeInfo(tx); err != nil {
			return err
		}

		if err := o.OrderLogs[0].Save(tx, o.ID); err != nil {
			return err
		}

		return r.SetToBusy(tx)
	})
}

// TriggerDriverCancel 订单取消
func (o *RepairRescueOrder) TriggerDriverCancel(d *Driver, r *Repair, fsm *RepairRescueFSM, c *config.OtherConf, res *Rescue, args *proto_order.OrderCancelArgs) error {
	times := d.dailyAvailCancel(c.DriverCancelTimes)
	if times == 0 {
		return ErrCancelNoTimes
	}

	fee, g, dis, err := o.ComputeMileFee(r, c)
	if err != nil {
		return err
	}

	if fee > 0 {
		return o.cancelWithFee(d, r, fsm, args.GetProcessBase().GetLocation(), fee, g, dis)
	}
	return o.cancelNoFee(d, r, fsm, res, args.GetProcessBase().GetLocation())
}

func (o *RepairRescueOrder) cancelWithFee(d *Driver, r *Repair, fsm *RepairRescueFSM, location *proto_order.Location, fee int64, g *Goods, dis float32) error {
	if err := fsm.Event(d.CancelFee()); err != nil {
		return err
	}
	o.Amount = fee

	o.modify(fsm.Current(), location, time.Now())

	o.OrderLogs = o.newLog(r.ID)
	oi := &OrderInfo{
		OrderInfoSchema: &OrderInfoSchema{
			GoodsID:  g.ID,
			Price:    g.Price,
			Name:     "里程费",
			Amount:   dis,
			InfoType: o.Type,
		},
	}
	o.OrderInfos = &OrderInfos{oi}
	return o.updateCancelWithFee(r)
}

func (o *RepairRescueOrder) cancelNoFee(d *Driver, r *Repair, fsm *RepairRescueFSM, res *Rescue, location *proto_order.Location) error {
	cf := func(tx *sql.Tx) error {
		return res.Cancel(tx)
	}
	return o.cancel(d.Cancel(), d.ID, r, fsm, location, cf)
}

// TriggerRepairCancel 订单取消
func (o *RepairRescueOrder) TriggerRepairCancel(r *Repair, fsm *RepairRescueFSM, c *config.OtherConf, args *proto_order.OrderCancelArgs) error {
	times := r.dailyAvailCancel(c.RepairCancelTimes)
	if times == 0 {
		return ErrCancelNoTimes
	}

	cf := func(tx *sql.Tx) error {
		r, err := LoadRescue(o.ID)
		if err != nil {
			return err
		}

		return r.Cancel(tx)
	}
	return o.cancel(r.Cancel(), r.ID, r, fsm, args.GetProcessBase().GetLocation(), cf)
}

func (o *RepairRescueOrder) cancel(event string, operator int64, r *Repair, fsm *RepairRescueFSM, location *proto_order.Location, cfs ...func(tx *sql.Tx) error) error {
	if err := fsm.Event(event); err != nil {
		return err
	}

	o.modify(fsm.Current(), location, time.Now())

	o.OrderLogs = o.newLog(operator)
	return o.updateCancel(r, cfs...)
}

func (o *RepairRescueOrder) updateCancelWithFee(r *Repair) error {
	return db.TxExec(func(tx *sql.Tx) error {
		if err := o.updatePriceChange(tx); err != nil {
			return err
		}

		// 技工设为接单中
		return r.SetToWork(tx)
	})
}

// TriggerRejectBy P2P订单拒绝接单
func (o *RepairRescueOrder) TriggerRejectBy(r *Repair, fsm *RepairRescueFSM, args *proto_order.OrderCancelArgs) error {
	return o.cancel(r.Cancel(), r.ID, r, fsm, args.GetProcessBase().GetLocation())
}

func (o *RepairRescueOrder) updateCancel(r *Repair, cfs ...func(tx *sql.Tx) error) error {
	return db.TxExec(func(tx *sql.Tx) error {
		if len(cfs) > 0 {
			for _, cf := range cfs {
				if err := cf(tx); err != nil {
					return err
				}
			}
		}

		if err := o.updateState(tx); err != nil {
			return err
		}

		// 技工设为接单中
		return r.SetToWork(tx)
	})
}

// ComputeMileFee 计算里程费
func (o *RepairRescueOrder) ComputeMileFee(r *Repair, c *config.OtherConf) (int64, *Goods, float32, error) {
	o.setRepair(r)
	g, err := LoadGoods(o.RepairOrgID)
	if err != nil {
		return 0, nil, 0, err
	}

	dis, err := o.computeDistance(r)
	if err != nil {
		return 0, nil, 0, err
	}

	if c.NeedPayFeeMiles > dis {
		return 0, nil, 0, nil
	}

	m := Money(float32(g.Price) * dis)
	return m.getFen(), g, dis, nil
}

func (o *RepairRescueOrder) updatePriceChange(tx *sql.Tx) error {
	if err := o.updatePrice(tx); err != nil {
		return err
	}

	// 保存订单日志
	return o.saveInfoAndLog(tx)
}

func (o *RepairRescueOrder) getCalculator(e, af int64) *RepairCalculator {
	return NewRepairCalculator(o.OrderInfos, e, af)
}

func (o *RepairRescueOrder) setCustomerPrice(e, af int64) {
	o.Extra = e
	o.AccessoriesFee = af
}

func (o *RepairRescueOrder) setRepair(r *Repair) {
	o.RepairID = r.ID
	o.RepairOrgID = r.OrgID
}

func (o *RepairRescueOrder) dispatcheTo(rs ...*RescueRepair) {
	sc, err := config.LoadSubsidyConf()
	if err != nil {
		//TODO 写日志
		return
	}

	for _, r := range rs {
		go r.recive(o, sc.Get())
	}
}

func (o *RepairRescueOrder) isCancel() bool {
	return cst.OrderStateCanceled.Equals(o.State) || cst.OrderStateCancelPay.Equals(o.State)
}

func (o *RepairRescueOrder) updateGrabeInfo(tx *sql.Tx) error {
	sqlStr := "UPDATE orders SET repair_id = ?, org_id = ?, status = ?, distance= ? WHERE id = ?"
	if _, err := db.UpdateTx(tx, sqlStr, o.RepairID, o.RepairOrgID, o.State, o.Distance, o.ID); err != nil {
		return err
	}
	return nil
}

func (o *RepairRescueOrder) computeDistance(r *Repair) (float32, error) {
	var dis float32
	sqlStr := `
	SELECT ROUND(SQRT(POW((? - ol.lat) * 111.15, 2) + POW((? - ol.lng) * 111.15, 2)), 2) AS distance
	FROM
		order_log ol
	WHERE
		order_id = ? AND status = 3
	LIMIT 1
	`
	err := db.DbSource.QueryRow(sqlStr, r.Location.lat, r.Location.lng, o.ID).Scan(&dis)
	return dis, err
}

// CloseRescueOrder 关闭救援订单
func CloseRescueOrder(driverID int64) error {
	sqlStr := `
	UPDATE orders SET status = 12
        WHERE user_id = ?
              AND (status = 0 AND place_order_type IN (2, 3)
              OR (order_type IN (1, 5, 13, 14, 15) AND (status = 5
                                                        OR (status = 6 AND EXISTS (SELECT 1
                                                                                  FROM
                                                                                      pay_others
                                                                                  WHERE
                                                                                      status = 3
                                                                                      AND id = (SELECT
                                                                                                    MAX(id)
                                                                                                FROM
                                                                                                    pay_others
                                                                                                WHERE
                                                                                                    apply_user_id = ?))))))
	`
	_, err := db.Update(sqlStr, driverID, driverID)
	if err != nil {
		return fmt.Errorf("autoClose order error: %v", err)
	}
	return nil
}
