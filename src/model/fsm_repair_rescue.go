package model

import (
	"ocenter/src/cst"
	"ocenter/src/model/config"

	"strconv"

	"github.com/looplab/fsm"
)

// RepairRescueFSM 维修
type RepairRescueFSM struct {
	*fsm.FSM
	Order  *RepairRescueOrder
	Driver *Driver
	Repair *Repair
	Config *config.OtherConf
}

// NewRepairRescueFSM 已有初始状态的加油订单状态机对象
func NewRepairRescueFSM(o *RepairRescueOrder) *RepairRescueFSM {
	r := &RepairRescueFSM{Order: o}
	r.initFSM()
	return r
}

func (r *RepairRescueFSM) initFSM() {
	r.FSM = fsm.NewFSM(
		strconv.Itoa(int(r.Order.State)),
		fsm.Events{
			// 创建事件	未创建 --> 创建
			{
				Name: cst.OrderOpEventPlace.String(),
				Src:  []string{cst.OrderStateNone.String()},
				Dst:  cst.OrderStateCreated.String(),
			},

			// 抢单事件	创建 --> 出发中
			{
				Name: cst.OrderOpEventRun.String(),
				Src:  []string{cst.OrderStateCreated.String()},
				Dst:  cst.OrderStateRuning.String(),
			},

			// 处理事件	开始修理 --> 完成
			{
				Name: cst.OrderOpEventFinish.String(),
				Src:  []string{cst.OrderStateRuning.String()},
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
			// 创建 --> 取消 | 出发中 --> 取消
			{
				Name: cst.OrderOpEventCancel.String(),
				Src: []string{
					cst.OrderStateCreated.String(),
					cst.OrderStateRuning.String(),
				},
				Dst: cst.OrderStateCanceled.String(),
			},
		},
		fsm.Callbacks{
			// before_create
			cst.FSMBefore + cst.OrderOpEventPlace.String(): r.beforeCreate,

			// before_run
			cst.FSMBefore + cst.OrderOpEventRun.String(): r.beforeRun,
		},
	)
}

// Event 委托事件
func (r *RepairRescueFSM) Event(e string) error {
	return r.FSM.Event(e)
}

// Current 当前状态
func (r *RepairRescueFSM) Current() string {
	return r.FSM.Current()
}

func (r *RepairRescueFSM) beforeCreate(e *fsm.Event) {
	if err := r.Driver.CanRescue(r.Config); err != nil {
		e.Err = err
		return
	}

	if r.Order.RepairID != 0 {
		e.Err = r.Repair.CheckStateAvailable()
	}
}

func (r *RepairRescueFSM) beforeRun(e *fsm.Event) {
	e.Err = r.Repair.CheckStateAvailable()
}
