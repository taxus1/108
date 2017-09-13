package model

import (
	"database/sql"
	"fmt"
	"ocenter/src/db"
	"ocenter/src/model/config"

	"time"
)

// 抢单状态
// 1创建订单资源
// 2已抢
// 3取消
const (
	statusCreate = 1
	statusGrabed = 2
	statusCancel = 3
)

var cols = []string{
	"order_id",
	"driver_id",
	"repair_id",
	"status",
	"lat",
	"lng",
	"address",
	"create_time",
	"update_time",
}

// RescueSchema 救援订单池
type RescueSchema struct {
	ID int64

	//订单ID
	OrderID int64

	//司机ID
	DriverID int64

	//技工ID
	RepairID int64

	//抢单状态
	State int32

	Lat float32

	Lng float32

	Address string

	CreateTime time.Time

	UpdateTime time.Time

	//抢单超时时间
	Timeout int32
}

// Rescue 救援订单池
type Rescue struct {
	*RescueSchema
}

// NewRescue 新建救援
func NewRescue(o *RepairRescueOrder) *Rescue {
	return &Rescue{
		RescueSchema: &RescueSchema{
			OrderID:    o.ID,
			DriverID:   o.DriverID,
			RepairID:   o.RepairID,
			State:      statusCreate,
			Lat:        o.Lat,
			Lng:        o.Lng,
			Address:    o.Address,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
	}
}

// LoadRescue 加载订单救援信息
func LoadRescue(orderID int64) (*Rescue, error) {
	sqlStr := `SELECT id, driver_id, status, create_time FROM rescues WHERE order_id = ?`
	r := &Rescue{RescueSchema: &RescueSchema{OrderID: orderID}}
	err := db.DbSource.QueryRow(sqlStr, orderID).Scan(
		&r.ID,
		&r.DriverID,
		&r.State,
		&r.CreateTime,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("load rescue error: %v", err)
	}
	return r, nil
}

// CheckGrabed 检查抢单
func (r *Rescue) CheckGrabed(c *config.GrabConf) error {
	switch {
	case r.State == 2:
		return ErrRescueGrabed
	case r.State == 3:
		return ErrRescueCanceled
	case time.Now().Unix()-r.CreateTime.Unix() > int64(c.WaiteTime*60):
		return ErrRescueExpired
	}
	return nil
}

// Grab 抢单
func (r *Rescue) Grab(tx *sql.Tx, repairID int64) error {
	sqlStr := `
		UPDATE rescues
		SET status = 2, repair_id = ?
		WHERE id = ? AND status = 1
	`
	n, err := db.UpdateTx(tx, sqlStr, repairID, r.ID)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrRescueExpired
	}
	return nil
}

// Cancel 救援取消
func (r *Rescue) Cancel(tx *sql.Tx) error {
	sqlStr := `UPDATE rescues SET status = 3 WHERE id = ?`
	if _, err := db.UpdateTx(tx, sqlStr, r.ID); err != nil {
		return err
	}
	return nil
}

func (r *Rescue) insert(tx *sql.Tx) error {
	sqlStr := db.BuildInsert("rescues", cols, 1)
	lastID, err := db.SaveTx(tx, sqlStr, r.values())
	if err != nil {
		return err
	}

	r.ID = lastID
	return nil
}

// values 获取对象中的值进行SQL设值
func (r *Rescue) values() []interface{} {
	return []interface{}{
		r.OrderID,
		r.DriverID,
		r.RepairID,
		r.State,
		r.Lat,
		r.Lng,
		r.Address,
		r.CreateTime,
		r.UpdateTime,
	}
}
