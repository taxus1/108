package model

import (
	"ocenter/src/model/config"
	"strconv"
)

// Calculator 计算器
type Calculator interface {
	Calculate() (int64, error)
}

// Price 价格
type Price struct {
	// 材料
	material *OrderInfos
}

// NewPrice 新建价格对象
func NewPrice(m *OrderInfos) *Price {
	return &Price{m}
}

func (p *Price) materialPriceFen() int64 {
	return p.materialPrice().getFen()
}

func (p *Price) materialPrice() *Money {
	return p.material.MaterialMoney()
}

func (p *Price) checkPrice(price int64) error {
	if price > 0 {
		return nil
	}
	return ErrOrderAmount
}

// RepairCalculator 维修金额
type RepairCalculator struct {
	*Price

	//里程费
	extra int64

	//附加费
	accessoriesFee int64
}

// NewRepairCalculator 创建维修金额对象
func NewRepairCalculator(m *OrderInfos, e, af int64) *RepairCalculator {
	return &RepairCalculator{NewPrice(m), e, af}
}

// Calculate 获取金额
func (r *RepairCalculator) Calculate() (int64, error) {
	p := r.materialPriceFen() + r.customerPrice()
	return p, r.checkPrice(p)
}

// customerPrice 用户自己输入的金额，附加费， 辅料费等
func (r *RepairCalculator) customerPrice() int64 {
	return r.extra + r.accessoriesFee
}

// OilCalculator 油品类订单金额计算
type OilCalculator struct {
	*Price
}

// NewOilCalculator 创建油品金额对象
func NewOilCalculator(m *OrderInfos) *OilCalculator {
	return &OilCalculator{Price: NewPrice(m)}
}

// Calculate 获取金额
func (o *OilCalculator) Calculate() (int64, error) {
	p := o.materialPriceFen()
	return p, o.checkPrice(p)
}

// InsuranceResuceCalculator 保险救援订单金额计算器
type InsuranceResuceCalculator struct {
	Org                 *Org
	InsuranceRescueRule *config.InsuranceRescueRule
}

// NewInsuranceResuceCalculator 保险救援订单金额计算对象
func NewInsuranceResuceCalculator(o *Org, irr *config.InsuranceRescueRule) *InsuranceResuceCalculator {
	return &InsuranceResuceCalculator{Org: o, InsuranceRescueRule: irr}
}

// Calculate 获取金额
func (o *InsuranceResuceCalculator) Calculate() int64 {
	if !o.InsuranceRescueRule.HasActivity() || !o.Org.Deposited() {
		return o.defaultPrice()
	}

	if o.InsuranceRescueRule.IsLimitTimeActivity() {
		if !o.InsuranceRescueRule.Expire(o.Org.ConfirmAt) {
			return o.InsuranceRescueRule.Amount
		}
		return o.defaultPrice()
	}

	if !o.InsuranceRescueRule.Overrun(o.Org.UnhandInsuranceOrders) {
		return o.InsuranceRescueRule.Amount
	}
	return o.defaultPrice()
}

func (o *InsuranceResuceCalculator) defaultPrice() int64 {
	r, _ := strconv.Atoi(o.InsuranceRescueRule.Content)
	return int64(r)
}
