package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"ocenter/src/cst"
	"ocenter/src/db"
	"ocenter/src/model/config"
	"ocenter/src/model/message"
	proto_order "ocenter/src/proto_order"

	"time"
)

// userModel 用户业务模型
type userModel Model

// User 用户公共类
type User struct {

	//ID
	ID int64

	//用户手机号
	Phone string

	// 用户名
	Name string

	//集团/维修厂ID
	OrgID int64

	//进行中的订单数
	ProcessOrders int32

	// 剩余救援次数
	RescueTimes int32

	Location *Location
}

// Cancel 取消
func (u *User) Cancel() string {
	return cst.OrderOpEventCancel.String()
}

// dailyAvailCancel 每日剩余可取消次数
func (u *User) dailyAvailCancel(times int32) int32 {
	if u.RescueTimes < 0 {
		return 0
	} else if u.RescueTimes > 0 {
		return u.RescueTimes
	}

	u.RescueTimes = u.computeAvailCancel(times)
	return u.RescueTimes
}

func (u *User) computeAvailCancel(times int32) int32 {
	canceled, err := u.countDailyCancelOrders()
	if err != nil {
		return -1
	}

	if times > canceled {
		return times - canceled
	}
	return -1
}

func (u *User) countDailyCancelOrders() (int32, error) {
	var count int32
	sqlStr := `SELECT
			            COUNT(1)
			        FROM
			            orders
			        WHERE
			            status = 9 AND user_id = ?
			            AND DATE(create_time) = ?
						`

	err := db.DbSource.QueryRow(sqlStr, u.ID, time.Now().Format("2006-01-02")).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}

// Repair 技工，由用户组成
type Repair struct {
	//用户类
	*User

	//技工用户状态
	State int32

	// 技工所在的组织
	Org *Org
}

// LoadRepair 根据ID获取技工
func LoadRepair(id int64) (r *Repair, err error) {
	sqlStr := "SELECT  IFNULL(u.status, 0) AS status, u.org_id, COUNT(o.id) AS orders, u.phone, uwp.lat, uwp.lng, uwp.address" +
		" FROM v_user_worker u LEFT JOIN `orders` o ON u.id = o.repair_id AND o.status < 5 LEFT JOIN user_worker_position uwp ON u.id = uwp.user_id" +
		" WHERE u.id = ?"
	r = &Repair{User: &User{ID: id, Location: &Location{}}}

	err = db.DbSource.QueryRow(sqlStr, id).Scan(
		&r.State,
		&r.OrgID,
		&r.ProcessOrders,
		&r.Phone,
		&r.Location.lat,
		&r.Location.lng,
		&r.Location.address,
	)
	return r, err
}

// SetToBusy 技工置为繁忙状态
func (r *Repair) SetToBusy(tx *sql.Tx) error {
	r.State = 7
	return r.updateState(tx)
}

// SetToWork 技工置为工作状态
func (r *Repair) SetToWork(tx *sql.Tx) error {
	r.State = 5
	return r.updateState(tx)
}

func (r *Repair) updateState(tx *sql.Tx) error {
	sqlStr := "UPDATE user SET status = ? WHERE id = ?"
	_, err := db.UpdateTx(tx, sqlStr, r.State, r.ID)
	if err != nil {
		return fmt.Errorf("update repair state error: %v", err)
	}
	return nil
}

// CheckStateAvailable 技工可用状态检测
func (r *Repair) CheckStateAvailable() error {
	switch r.State {
	case cst.RepairRest:
		return ErrRepairStateRest
	case cst.RepairBusy:
		return ErrRepairStateBusy
	case cst.RepairWork:
		return nil
	default:
		return ErrRepairStateNotWork
	}
}

// PlaceOrder 下单
func (r *Repair) PlaceOrder() string {
	return cst.OrderOpEventPlace.String()
}

// Bargain 议价
func (r *Repair) Bargain() string {
	return cst.OrderOpEventBargain.String()
}

// Run 出发
func (r *Repair) Run() string {
	return cst.OrderOpEventRun.String()
}

// Grab 抢单
func (r *Repair) Grab() string {
	return r.Run()
}

// Accept 接受订单
func (r *Repair) Accept() string {
	return r.Run()
}

// Hand 修理
func (r *Repair) Hand() string {
	return cst.OrderOpEventHand.String()
}

// Finish 完成修理
func (r *Repair) Finish() string {
	return cst.OrderOpEventFinish.String()
}

// Pay 支付保险救援订单
func (r *Repair) Pay() string {
	return cst.OrderOpEventPay.String()
}

// CheckOrgInsuranceOrder 检查技工所在维修厂未处理保险救援订单
func (r *Repair) CheckOrgInsuranceOrder() error {
	o, err := LoadOrg(r.OrgID)
	if err != nil {
		return err
	}
	return o.CheckInsuranceOrders()
}

// GetOrg 获取技工所在维修厂
func (r *Repair) GetOrg() (*Org, error) {
	o, err := LoadOrg(r.OrgID)
	if err != nil {
		return nil, err
	}

	r.Org = o
	return o, nil
}

// RescueRepair 救援维修工
type RescueRepair struct {
	ID             int64
	Distance       float32
	Address        string
	Realname       string
	OrderBackScore int32
	OrderBackTimes int32
	Phone          string
	HeaderImg      string
	totalOrder     int32
}

// FindRepairor 查询满足条件的技工
func FindRepairor(d *Driver, l *proto_order.Location, troubles *orderInfoTypes) ([]*RescueRepair, error) {
	c, err := config.LoadGrabConf()
	if err != nil {
		//TODO 写日志
		return nil, err
	}

	buf := bytes.Buffer{}
	buf.WriteString(`
		 SELECT * FROM
		            (SELECT
		                u.id AS userId,
		                ROUND(SQRT(POW((? - urp.lat) * 111.15, 2) + POW((? - urp.lng) * 111.15, 2)), 2) AS distance,
		                urp.address,
		                u.realname,
		                u.order_back_score AS orderBackScore,
		                u.order_back_times AS orderBackTimes,
		                u.phone,
		                u.header_img AS headerImg,
		                u.order_total As totalOrder
		            FROM v_user_worker u
		                RIGHT JOIN truck_org_worker tow ON tow.worker_id = u.id
		                LEFT JOIN org o ON o.id = tow.org_id
		                LEFT JOIN user_worker_position urp ON u.id = urp.user_id
		                LEFT JOIN org_info oi ON oi.org_id = o.id
		            WHERE u.status = 5
		                AND u.user_type IN (4, 5)
		                AND ROUND(SQRT(POW((? - urp.lat) * 111.15, 2) + POW((? - urp.lng) * 111.15, 2)), 2) <= ?
		                AND o.deleted = 0
		                AND oi.is_pact = 1
		                AND urp.lat != 0
		                AND urp.lng != 0
	`)
	if len(*troubles) > 0 {
		d.CurrentTruck.TroubleCondition(&buf, troubles)
	}
	buf.WriteString(`
	        ) tmp
	        WHERE tmp.distance >= ?
	        ORDER BY distance
	        LIMIT ?
	`)

	var rrs []*RescueRepair
	f := func(rs *sql.Rows) error {
		for rs.Next() {
			rr := new(RescueRepair)
			if err := rs.Scan(rr.fields()); err != nil {
				return err
			}
			rrs = append(rrs, rr)
		}
		return nil
	}
	return rrs, db.QueryMore(
		buf.String(),
		f,
		l.Lat,
		l.Lng,
		l.Lat,
		l.Lng,
		c.MaxRepairMiles,
		c.MinRepairMiles,
		c.GrabNum,
	)
}

func (rr *RescueRepair) fields() []interface{} {
	return []interface{}{
		&rr.ID,
		&rr.Distance,
		&rr.Address,
		&rr.Realname,
		&rr.OrderBackScore,
		&rr.OrderBackTimes,
		&rr.Phone,
		&rr.HeaderImg,
		&rr.totalOrder,
	}
}

func (rr *RescueRepair) recive(o *RepairRescueOrder, subsidy string) {
	m := &message.RescueOrderMsg{
		OrderID:        o.ID,
		Distance:       rr.Distance,
		Subsidy:        subsidy,
		PlaceOrderType: 2,
	}
	message.Sender.Send(rr.ID, m)
}

// Driver 司机，由用户组成
type Driver struct {
	//用户类
	*User

	CurrentTruck *Truck
}

// LoadDriver id获取用户信息
func LoadDriver(id int64) (*Driver, error) {
	sqlStr := "SELECT u.realname, u.org_id, COUNT(o.id) AS orders, u.phone, uwp.lat, uwp.lng, uwp.address" +
		" FROM v_user_driver u LEFT JOIN `orders` o ON u.id = o.user_id AND o.status <= 5 LEFT JOIN user_driver_position uwp ON u.id = uwp.user_id" +
		" WHERE u.id = ? "

	d := &Driver{User: &User{ID: id, Location: &Location{}}}

	err := db.DbSource.QueryRow(sqlStr, id).Scan(
		&d.Name,
		&d.OrgID,
		&d.ProcessOrders,
		&d.Phone,
		&d.Location.lat,
		&d.Location.lng,
		&d.Location.address,
	)
	if err != nil {
		return nil, err
	}
	return d, err
}

// LoadCurrentTruck 司机当前车辆
func (d *Driver) LoadCurrentTruck() error {
	t, err := LoadCurrentTruck(d.ID)
	if err != nil {
		return err
	}
	d.CurrentTruck = t
	return nil
}

// GetCurrentTruck 获取司机当前车辆
func (d *Driver) GetCurrentTruck() *Truck {
	if d != nil {
		return d.CurrentTruck
	}
	return nil
}

// CanRepair 司机能否进行维修业务
func (d *Driver) CanRepair(r *Repair) error {
	// 检测司机未完成订单
	if err := d.CheckProcessOrder(); err != nil {
		return err
	}

	// 检测司机车辆信息
	if err := d.CheckTruckRepairable(); err != nil {
		return err
	}

	// 检查下单司机和技工的手机号不能相同
	return d.PhoneEqualRepairPhone(r)
}

// CanFuel 司机能否进行加油业务
func (d *Driver) CanFuel() error {

	// 检测司机车辆信息
	return d.CheckTruckFuelable()
}

// CanRescue 司机能否进行救援业务
func (d *Driver) CanRescue(c *config.OtherConf) error {
	// 检测司机未完成订单
	if err := d.CheckProcessOrder(); err != nil {
		return err
	}

	// 检测司机车辆信息
	if err := d.CheckTruckRescueable(); err != nil {
		return err
	}

	if t := d.dailyAvailCancel(c.DriverCancelTimes); t <= 0 {
		return ErrCancelTooMany
	}
	return nil
}

// CanInsuranceRescue 司机能否进行保险救援业务
func (d *Driver) CanInsuranceRescue() error {
	// 检测司机未完成订单
	if err := d.CheckProcessOrder(); err != nil {
		return err
	}

	// 检测司机车辆信息
	return d.CheckTruckRescueable()
}

// CheckProcessOrder 司机是否有正在进行中的订单
func (d *Driver) CheckProcessOrder() error {
	if d.ProcessOrders > 0 {
		return ErrDriverHaveProcessOrder
	}

	return nil
}

//CheckTruckRepairable 司机检测当前车辆是否能修理
func (d *Driver) CheckTruckRepairable() error {
	if err := d.LoadCurrentTruck(); err != nil {
		return err
	}
	return d.CurrentTruck.CheckRepaireable()
}

//CheckTruckFuelable 司机检测当前车辆是否能加油
func (d *Driver) CheckTruckFuelable() error {
	if err := d.LoadCurrentTruck(); err != nil {
		return err
	}
	return d.CurrentTruck.CheckFuelable()
}

//CheckTruckRescueable 司机检测当前车辆是否能发起救援
func (d *Driver) CheckTruckRescueable() error {
	if err := d.LoadCurrentTruck(); err != nil {
		return err
	}
	return d.CurrentTruck.CheckRescueable()
}

// PhoneEqualRepairPhone 检查司机和技工手机号
func (d *Driver) PhoneEqualRepairPhone(r *Repair) error {
	if d.Phone == r.Phone {
		return ErrDriverPhoneEqualsRepairPhone
	}

	return nil
}

// PlaceOrder 下单
func (d *Driver) PlaceOrder() string {
	return cst.OrderOpEventPlace.String()
}

// Confirm 确认
func (d *Driver) Confirm() string {
	return cst.OrderOpEventConfirm.String()
}

// OtherPay 代付
func (d *Driver) OtherPay() string {
	return cst.OrderOpEventOtherPay.String()
}

// Pay 支付
func (d *Driver) Pay() string {
	return cst.OrderOpEventPay.String()
}

// Assess 支付
func (d *Driver) Assess() string {
	return cst.OrderOpEventAssess.String()
}

// Close 关闭
func (d *Driver) Close() string {
	return cst.OrderOpEventClose.String()
}

// CancelFee 带里程费取消
func (d *Driver) CancelFee() string {
	return cst.OrderOpEventCancelFee.String()
}
