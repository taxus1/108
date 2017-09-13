package model

import (
	"ocenter/src/cst"

	"strconv"

	"github.com/looplab/fsm"
)

// OilFSM 加油业务状态机
type OilFSM struct {
	*fsm.FSM
	Order  *OilOrder
	Driver *Driver
}

// NewOilFSM 已有初始状态的加油订单状态机对象
func NewOilFSM(o *OilOrder) *OilFSM {
	r := &OilFSM{Order: o}
	r.initFSM()
	return r
}

func (r *OilFSM) initFSM() {
	r.FSM = fsm.NewFSM(
		strconv.Itoa(int(r.Order.State)),
		fsm.Events{
			// 创建事件	未创建 --> 完成
			{
				Name: cst.OrderOpEventPlace.String(),
				Src:  []string{cst.OrderStateNone.String()},
				Dst:  cst.OrderStateFinished.String(),
			},

			// 代付事件	完成 --> 代付
			{
				Name: cst.OrderOpEventOtherPay.String(),
				Src:  []string{cst.OrderStateFinished.String()},
				Dst:  cst.OrderStateOtherPayed.String(),
			},

			// 支付事件	完成 --> 支付 | 代付 --> 支付
			{
				Name: cst.OrderOpEventPay.String(),
				Src: []string{
					cst.OrderStateFinished.String(),
					cst.OrderStateOtherPayed.String(),
				},
				Dst: cst.OrderStatePayed.String(),
			},

			// 评价事件 支付 --> 评价
			{
				Name: cst.OrderOpEventAssess.String(),
				Src:  []string{cst.OrderStatePayed.String()},
				Dst:  cst.OrderStateAssessed.String(),
			},

			// 取消事件
			// 创建 --> 取消
			{
				Name: cst.OrderOpEventClose.String(),
				Src:  []string{cst.OrderStateFinished.String()},
				Dst:  cst.OrderStateClosed.String(),
			},
		},
		fsm.Callbacks{
			// before_create
			cst.FSMBefore + cst.OrderOpEventPlace.String(): r.beforeCreate,
		},
	)
}

// Event 委托事件
func (r *OilFSM) Event(e string) error {
	return r.FSM.Event(e)
}

// Current 当前状态
func (r *OilFSM) Current() string {
	return r.FSM.Current()
}

func (r *OilFSM) beforeCreate(e *fsm.Event) {
	e.Err = r.Driver.CanFuel()
}
