package model

import (
	"ocenter/src/db"
	"time"
)

// 服务费类型，1 按数量收取，2 按总金额收取
const (
	OilChargeCount  = 1
	OilChargeAmount = 2
)

// Org 维修厂
type Org struct {
	ID                    int64
	State                 int32
	Manager               int64
	Trailer               int32
	Metaller              int32
	Electrician           int32
	Mechanic              int32
	Insurance             string
	ConfirmAt             time.Time
	OilChargeType         int32
	OilChargeAmount       int64
	UnhandInsuranceOrders int32
}

// LoadOrg 加载维修厂
func LoadOrg(id int64) (*Org, error) {
	sqlStr := `
	SELECT
		o.oil_charge_type,
		o.oil_charge_amount,
		ore.status,
		ore.repair_id,
		ore.trailer,
		ore.metaller,
		ore.electrician,
		ore.mechanic,
		ore.insurance,
		ore.update_time,
		tmp.unhands
	FROM
		org o
				LEFT JOIN
		org_rescue ore ON o.id = ore.org_id,
		(SELECT
				COUNT(1) unhands
		FROM
				orders
		WHERE
				org_id = ? AND status = 5
						AND order_type = 19) tmp
	WHERE
		o.id = ?
  `
	o := &Org{ID: id}
	err := db.DbSource.QueryRow(sqlStr, id).Scan(
		&o.OilChargeType,
		&o.OilChargeAmount,
		&o.State,
		&o.Manager,
		&o.Trailer,
		&o.Metaller,
		&o.Electrician,
		&o.Mechanic,
		&o.Insurance,
		&o.ConfirmAt,
		&o.UnhandInsuranceOrders,
	)
	return o, err
}

// CheckInsuranceOrders 检查未处理的保险救援订单
func (o *Org) CheckInsuranceOrders() error {
	if o.UnhandInsuranceOrders != 0 {
		return ErrOrgHasInsuranceOrder
	}
	return nil
}

// Deposited 是否缴纳过保证金
func (o *Org) Deposited() bool {
	return o.State == 1
}

// ServiceFee 服务费
func (o *Org) ServiceFee(count float32, amount int64) int64 {
	switch o.OilChargeType {
	case OilChargeCount:
		m := Money(o.OilChargeAmount)
		m.multiply(count)
		return m.getFen()
	case OilChargeAmount:
		return o.OilChargeAmount * amount
	default:
		return 0
	}
}
