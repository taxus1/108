// Package cst 状态机常量
package cst

import (
	"strconv"
)

//EventInt 订单流程事件常量
type EventInt int

//S 数字常量转字符类型
func (i EventInt) String() string {
	return "event_" + strconv.Itoa(int(i))
}

//订单事件
const (
	OrderOpEventNone      EventInt = -1
	OrderOpEventPlace     EventInt = 0
	OrderOpEventBargain   EventInt = 1
	OrderOpEventConfirm   EventInt = 2
	OrderOpEventRun       EventInt = 3
	OrderOpEventHand      EventInt = 4
	OrderOpEventFinish    EventInt = 5
	OrderOpEventOtherPay  EventInt = 6
	OrderOpEventPay       EventInt = 7
	OrderOpEventAssess    EventInt = 8
	OrderOpEventCancel    EventInt = 9
	OrderOpEventClose     EventInt = 10
	OrderOpEventFail      EventInt = 11
	OrderOpEventCancelFee EventInt = 13
)

//OrderStateInt 订单状态常量
type OrderStateInt int

//S 数字常量转字符类型
func (i OrderStateInt) String() string {
	return strconv.Itoa(int(i))
}

// Equals 值是否相等
func (i OrderStateInt) Equals(s int32) bool {
	return int32(i) == s
}

//订单状态
const (
	OrderStateNone       OrderStateInt = -1
	OrderStateCreated    OrderStateInt = 0
	OrderStateBargain    OrderStateInt = 1
	OrderStateConfirmed  OrderStateInt = 2
	OrderStateRuning     OrderStateInt = 3
	OrderStateHanding    OrderStateInt = 4
	OrderStateFinished   OrderStateInt = 5
	OrderStateOtherPayed OrderStateInt = 6
	OrderStatePayed      OrderStateInt = 7
	OrderStateAssessed   OrderStateInt = 8
	OrderStateCanceled   OrderStateInt = 9
	OrderStateFailed     OrderStateInt = 11
	OrderStateClosed     OrderStateInt = 12
	OrderStateCancelPay  OrderStateInt = 13
)

//状态机默认事件
const (
	FSMBefore = "before_"
	FSMEnter  = "enter_"
	FSMLeave  = "leave_"
	FSMAfter  = "after_"
)
