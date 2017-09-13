package model

import (
	"ocenter/src/cst"

	"strconv"

	"github.com/looplab/fsm"
)

// InsuranceRescueFSM 保险救援状态机
type InsuranceRescueFSM struct {
	*fsm.FSM
	Order  *InsuranceRescueOrder
	Driver *Driver
	Repair *Repair
}

// NewInsuranceRescueFSM 已有初始状态的保险救援订单状态机对象
func NewInsuranceRescueFSM(o *InsuranceRescueOrder) *InsuranceRescueFSM {
	r := &InsuranceRescueFSM{Order: o}
	r.initFSM()
	return r
}

func (r *InsuranceRescueFSM) initFSM() {
	r.FSM = fsm.NewFSM(
		strconv.Itoa(int(r.Order.State)),
		fsm.Events{
			// 创建事件	未创建 --> 创建
			{
				Name: cst.OrderOpEventPlace.String(),
				Src:  []string{cst.OrderStateNone.String()},
				Dst:  cst.OrderStateFinished.String(),
			},

			// 支付事件	完成 --> 支付 | 代付 --> 支付
			{
				Name: cst.OrderOpEventPay.String(),
				Src:  []string{cst.OrderStateFinished.String()},
				Dst:  cst.OrderStatePayed.String(),
			},

			// 评价事件 支付 --> 评价
			{
				Name: cst.OrderOpEventAssess.String(),
				Src:  []string{cst.OrderStatePayed.String()},
				Dst:  cst.OrderStateAssessed.String(),
			},

			// 取消事件
			// 创建 --> 取消 | 出发中 --> 取消
			{
				Name: cst.OrderOpEventCancel.String(),
				Src:  []string{cst.OrderStateFinished.String()},
				Dst:  cst.OrderStateCanceled.String(),
			},
		},
		fsm.Callbacks{
			// before_create
			cst.FSMBefore + cst.OrderOpEventPlace.String(): r.beforeCreate,
		},
	)
}

// Event 委托事件
func (r *InsuranceRescueFSM) Event(e string) error {
	return r.FSM.Event(e)
}

// Current 当前状态
func (r *InsuranceRescueFSM) Current() string {
	return r.FSM.Current()
}

func (r *InsuranceRescueFSM) beforeCreate(e *fsm.Event) {
	// 检测司机未完成订单
	if err := r.Driver.CanInsuranceRescue(); err != nil {
		e.Err = err
		return
	}

	if err := r.Repair.CheckStateAvailable(); err != nil {
		e.Err = err
		return
	}

	e.Err = r.Repair.CheckOrgInsuranceOrder()
}
