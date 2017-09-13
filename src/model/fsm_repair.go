package model

import (
	"ocenter/src/cst"

	"strconv"

	"github.com/looplab/fsm"
)

// RepairFSM 维修
type RepairFSM struct {
	*fsm.FSM
	Order  *RepairOrder
	Driver *Driver
	Repair *Repair
}

// NewRepairFSM 已有初始状态的加油订单状态机对象
func NewRepairFSM(o *RepairOrder) *RepairFSM {
	f := &RepairFSM{Order: o}
	f.initFSM()
	return f
}

func (r *RepairFSM) initFSM() {
	r.FSM = fsm.NewFSM(
		strconv.Itoa(int(r.Order.State)),
		fsm.Events{
			// 创建事件	未创建 --> 创建
			{
				Name: cst.OrderOpEventPlace.String(),
				Src:  []string{cst.OrderStateNone.String()},
				Dst:  cst.OrderStateCreated.String(),
			},

			// 议价事件	创建 --> 议价
			{
				Name: cst.OrderOpEventBargain.String(),
				Src:  []string{cst.OrderStateCreated.String()},
				Dst:  cst.OrderStateBargain.String(),
			},

			// 确认事件	议价 --> 已确认
			{
				Name: cst.OrderOpEventConfirm.String(),
				Src:  []string{cst.OrderStateBargain.String()},
				Dst:  cst.OrderStateConfirmed.String(),
			},

			// 出发事件	确认 --> 出发中
			{
				Name: cst.OrderOpEventRun.String(),
				Src:  []string{cst.OrderStateConfirmed.String()},
				Dst:  cst.OrderStateRuning.String(),
			},

			// 处理事件	出发中 --> 开始修理
			{
				Name: cst.OrderOpEventHand.String(),
				Src:  []string{cst.OrderStateRuning.String()},
				Dst:  cst.OrderStateHanding.String(),
			},

			// 处理事件	开始修理 --> 完成
			{
				Name: cst.OrderOpEventFinish.String(),
				Src:  []string{cst.OrderStateHanding.String()},
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
			// 创建 --> 取消 | 议价 --> 取消 | 确认 --> 取消 | 出发中 --> 取消 | 维修中 --> 取消
			{
				Name: cst.OrderOpEventCancel.String(),
				Src: []string{
					cst.OrderStateCreated.String(),
					cst.OrderStateBargain.String(),
					cst.OrderStateConfirmed.String(),
					cst.OrderStateRuning.String(),
					cst.OrderStateHanding.String(),
				},
				Dst: cst.OrderStateCanceled.String(),
			},
		},
		fsm.Callbacks{
			// before_create
			cst.FSMBefore + cst.OrderOpEventPlace.String(): r.beforeCreate,
		},
	)
}

// Event 委托事件
func (r *RepairFSM) Event(e string) error {
	return r.FSM.Event(e)
}

// Current 当前状态
func (r *RepairFSM) Current() string {
	return r.FSM.Current()
}

func (r *RepairFSM) beforeCreate(e *fsm.Event) {
	// 检测司机未完成订单
	if err := r.Driver.CanRepair(r.Repair); err != nil {
		e.Err = err
		return
	}

	// 检测技工状态
	if err := r.Repair.CheckStateAvailable(); err != nil {
		e.Err = err
		return
	}

	// 检查下单司机和技工的手机号不能相同
	e.Err = r.Driver.PhoneEqualRepairPhone(r.Repair)
}
