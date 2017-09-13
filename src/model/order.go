package model

import (
	"database/sql"
	"fmt"
	"ocenter/src/cst"
	"ocenter/src/db"
	"ocenter/src/util"
	"strconv"
	"time"

	proto_order "ocenter/src/proto_order"
)

//OrderCols 订单表列名
var OrderCols = []string{
	"order_code",
	"user_id",
	"org_id",
	"order_type",
	"status",
	"amount",
	"lat",
	"lng",
	"remark",
	"area_id",
	"is_come",
	"repair_id",
	"repair_truck_id",
	"truck_id",
	"create_time",
	"finish_time",
	"order_address",
	"extra",
	"order_sub_type",
	"invoice_log_id",
	"distance",
	"deleted",
	"accessories_fee",
	"mileage",
	"coupon_id",
	"oil_gun",
	"order_time",
	"service_time",
	"place_order_type",
	"service_fee",
}

// orderModel 订单
type orderModel Model

// OrderSchema 订单表定义
type OrderSchema struct {
	//id
	ID int64

	//订单号
	Code string

	//所属用户
	DriverID int64

	//车辆ID
	DriverTruckID int64

	//技工ID
	RepairID int64

	//技工的车
	RepairTruckID int64

	//技工组织ID
	RepairOrgID int64

	//订单类型
	Type int32

	//订单状态
	State int32

	//订单总价格
	Amount int64

	//坐标lat
	Lat float32

	//坐标lng
	Lng float32

	//备注
	Remark string

	//所在地区ID
	AreaID int32

	//是否上门服务
	IsCome int32

	//订单下单时间
	CreateTime time.Time

	//订单完成时间
	FinishTime time.Time

	//订单中文位置
	Address string

	//附加费
	Extra int64

	//订单子类型
	SubType string

	//发票信息ID
	InvoiceLogID int64

	//距离
	Distance float32

	//技工删除
	Deleted int32

	//辅料费
	AccessoriesFee int64

	//里程费补贴
	Mileage int64

	//优惠券ID
	CouponID int64

	//油枪号
	OilGun string

	//下单方式
	PlaceOrderType int32

	//车牌号
	PlateNumber string

	//订单可提现时间
	RewardPeriod time.Time

	//车辆参保公司
	Insurance string

	//接单时长
	OrderTime int32

	//服务时长
	ServiceTime int32

	//服务费
	ServiceFee int64
}

//Order 订单对象
type Order struct {
	*OrderSchema

	OrderInfos *OrderInfos

	OrderLogs []*OrderLog

	SubType *orderInfoTypes
}

// NewOrder 订单超类对象
func NewOrder(t int32) *Order {
	return &Order{OrderSchema: &OrderSchema{
		Type:       t,
		Code:       util.GenCode(),
		CreateTime: time.Now(),
		FinishTime: time.Now(),
		State:      int32(cst.OrderStateNone),
	}}
}

func (o *Order) setDriver(d *Driver) {
	o.DriverID = d.ID
}

func (o *Order) setRepair(r *Repair) {
	o.RepairID = r.ID
	o.RepairOrgID = r.OrgID
}

// SetState 设置订单状态
func (o *Order) setState(c string) {
	i, _ := strconv.ParseInt(c, 10, 32)
	o.State = int32(i)
}

// SetLocation 设置订单地址
func (o *Order) setLocation(l *proto_order.Location) {
	if l == nil {
		return
	}
	o.Lat = l.Lat
	o.Lng = l.Lng
	o.Address = l.OrderAddress
	o.AreaID = l.AreaID
}

func (o *Order) setTruck(id int64) {
	o.DriverTruckID = id
}

func (o *Order) modify(state string, location *proto_order.Location, time time.Time) {
	o.setState(state)
	o.setLocation(location)
	o.FinishTime = time
}

func (o *Order) newInfo(infos []*proto_order.OrderInfo) *OrderInfos {
	return NewOrderInfos(infos, o.Type)
}

func (o *Order) newLog(operator int64) []*OrderLog {
	return []*OrderLog{NewOrderLog(o, operator)}
}

// 订单失效
func (o *Order) isInvalid() bool {
	return cst.OrderStateCanceled.Equals(o.State) || cst.OrderStateClosed.Equals(o.State)
}

func (o *Order) isPaid() bool {
	return cst.OrderStatePayed.Equals(o.State) || cst.OrderStateAssessed.Equals(o.State)
}

func (o *Order) setAmount(c Calculator) error {
	p, err := c.Calculate()
	if err != nil {
		return err
	}
	o.Amount = p
	return nil
}

func (o *Order) subTypeStr() string {
	if o.SubType == nil {
		return ""
	}
	return o.SubType.join(",")
}

// TriggerOtherPayBy 触发司机申请代付
func (o *Order) TriggerOtherPayBy(d *Driver, fsm EventCurrent, args *proto_order.OrderOtherPayArgs) error {
	if err := fsm.Event(d.OtherPay()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(d.ID)
	return o.updateOtherPay()
}

// TriggerPayBy 司机付款
func (o *Order) TriggerPayBy(d *Driver, fsm EventCurrent, args *proto_order.OrderPayArgs) error {
	if err := fsm.Event(d.Pay()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(d.ID)
	return o.updatePay()
}

// TriggerAssessBy 司机评价
func (o *Order) TriggerAssessBy(d *Driver, fsm EventCurrent, args *proto_order.OrderAssessArgs) error {
	if err := fsm.Event(d.Assess()); err != nil {
		return err
	}

	o.modify(fsm.Current(), args.GetProcessBase().GetLocation(), time.Now())
	o.OrderLogs = o.newLog(d.ID)
	return o.updateAccess()
}

func (o *Order) updateOtherPay() error {
	return o.updateBase()
}

//更新支付
func (o *Order) updatePay() error {
	return o.updateBase()
}

//更新评价
func (o *Order) updateAccess() error {
	return o.updateBase()
}

// saveRelative 保存订单，订单信息，订单日志
func (o *Order) saveRelative(fs ...func(tx *sql.Tx) error) error {
	return db.TxExec(func(tx *sql.Tx) error {
		// 保存订单
		if err := o.insert(tx); err != nil {
			return err
		}

		// 保存订单日志
		if err := o.saveInfoAndLog(tx); err != nil {
			return err
		}

		// 保存各种类型的订单自定义需要保存的信息
		if len(fs) == 0 {
			return nil
		}
		return fs[0](tx)
	})
}

func (o *Order) saveInfoAndLog(tx *sql.Tx) error {
	// 保存订单日志
	if err := o.OrderLogs[0].Save(tx, o.ID); err != nil {
		return err
	}

	// 保存订单信息
	return o.OrderInfos.BatchInsert(tx, o.ID, o.State)
}

// 更新状态和添加日志
func (o *Order) updateBase() error {
	return db.TxExec(func(tx *sql.Tx) error {
		return o.updateState(tx)
	})
}

func (o *Order) updateState(tx *sql.Tx) error {
	if err := o.update(tx); err != nil {
		return err
	}

	// 保存订单日志
	return o.OrderLogs[0].Save(tx, o.ID)
}

// Insert 插入订单
func (o *Order) insert(tx *sql.Tx) error {
	sqlStr := db.BuildInsert("orders", OrderCols, 1)
	lastID, err := db.SaveTx(tx, sqlStr, o.Values())
	if err != nil {
		return fmt.Errorf("insert order error: %s", err.Error())
	}

	o.ID = lastID
	return nil
}

// GetByID 根据ID获取订单
func (o *Order) load() error {
	query := "SELECT  id, order_code, user_id, org_id, order_type, status, repair_id," +
		" truck_id, create_time, finish_time, amount, extra, invoice_log_id, " +
		" distance, accessories_fee, mileage, coupon_id, oil_gun, place_order_type," +
		" plate_number, reward_period, insurance FROM `orders` WHERE id = ? LIMIT 1"
	fields := o.fields()
	return db.DbSource.QueryRow(query, o.ID).Scan(fields...)
}

// updatePrice 更新技工议价后的订单
func (o *Order) updatePrice(tx *sql.Tx) error {
	sqlStr := "UPDATE orders SET status = ?, amount = ?, finish_time = ?, extra = ?, accessories_fee = ? WHERE id = ?"
	_, err := db.UpdateTx(tx, sqlStr, o.State, o.Amount, o.FinishTime, o.Extra, o.AccessoriesFee, o.ID)
	if err != nil {
		return fmt.Errorf("update order price error: %v", err)
	}
	return nil
}

func (o *Order) update(tx *sql.Tx) error {
	sqlStr := "UPDATE orders SET status = ?, finish_time = ? WHERE id = ?"
	_, err := db.UpdateTx(tx, sqlStr, o.State, o.FinishTime, o.ID)
	if err != nil {
		return fmt.Errorf("update order state error: %v", err)
	}
	return nil
}

// Values 获取对象中的值进行SQL设值
func (o *Order) Values() []interface{} {
	return []interface{}{
		o.Code,
		o.DriverID,
		o.RepairOrgID,
		o.Type,
		o.State,
		o.Amount,
		o.Lat,
		o.Lng,
		o.Remark,
		o.AreaID,
		o.IsCome,
		o.RepairID,
		o.RepairTruckID,
		o.DriverTruckID,
		o.CreateTime,
		o.FinishTime,
		o.Address,
		o.Extra,
		o.subTypeStr(),
		o.InvoiceLogID,
		o.Distance,
		o.Deleted,
		o.AccessoriesFee,
		o.Mileage,
		o.CouponID,
		o.OilGun,
		o.OrderTime,
		o.ServiceTime,
		o.PlaceOrderType,
		o.ServiceFee,
	}
}

//Fields 获取订单字段
func (o *Order) fields() []interface{} {
	return []interface{}{
		&o.ID,
		&o.Code,
		&o.DriverID,
		&o.RepairOrgID,
		&o.Type,
		&o.State,
		&o.RepairID,
		&o.DriverTruckID,
		&o.CreateTime,
		&o.FinishTime,
		&o.Amount,
		&o.Extra,
		&o.InvoiceLogID,
		&o.Distance,
		&o.AccessoriesFee,
		&o.Mileage,
		&o.CouponID,
		&o.OilGun,
		&o.PlaceOrderType,
		&o.PlateNumber,
		&o.RewardPeriod,
		&o.Insurance,
	}
}
