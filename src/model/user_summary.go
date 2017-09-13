package model

import (
	"database/sql"
	"ocenter/src/db"
)

// userSummaryModel 用户统计模型
type userSummaryModel Model

// UserSummary 用户统计
type UserSummary struct {
	ID int64

	//用户ID
	UserID int64

	//接单时长
	OrderTime int32

	//接单次数
	OrderCount int32

	//服务时长
	ServiceTime int32

	//服务次数
	ServiceCount int32

	//付款订单数
	PayOrders int32
}

// UpdateOrderTime 更新技工响应订单的时间
func (m *userSummaryModel) UpdateOrderTime(tx *sql.Tx, time int32, repair int64) error {
	sqlStr := "UPDATE user_summaries SET order_times = order_times + ?, order_count = order_count + 1 WHERE user_id = ?"
	eff, err := db.UpdateTx(tx, sqlStr, []interface{}{time, repair})
	if err != nil {
		return err
	}

	if eff == 0 {
		us := &UserSummary{UserID: repair, OrderTime: time}
		return us.save(tx)
	}
	return nil
}

func (us *UserSummary) save(tx *sql.Tx) (err error) {
	cols := []string{"user_id", "order_times", "order_count", "service_times", "service_count"}
	sqlStr := db.BuildInsert("user_summaries", cols, len(cols))

	_, err = db.SaveTx(tx, sqlStr, us.values())
	return
}

func (us *UserSummary) values() []interface{} {
	return []interface{}{&us.UserID, &us.OrderTime, &us.OrderCount, &us.ServiceTime, &us.ServiceCount}
}
