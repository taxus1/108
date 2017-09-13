package model

import "errors"

// UserErr 用户相关错误
type UserErr error

var (

	// ErrRepairStateRest 当司机对技工下单，技工自己下单，技工状态为休息中返回该错误，
	ErrRepairStateRest UserErr = errors.New("技工休息中")

	// ErrRepairStateBusy 下单技工状态为繁忙，时返回该错误
	ErrRepairStateBusy UserErr = errors.New("技工繁忙，请稍后")

	// ErrRepairStateNotWork 司机对技工下单，技工自己下单，都需要技工在工作中的状态才能进行
	// 即状态为5时
	ErrRepairStateNotWork UserErr = errors.New("技工需要在上班状态才能接受订单")

	// ErrDriverPhoneEqualsRepairPhone 下单时司机和技工的手机号相同则返回该错误
	ErrDriverPhoneEqualsRepairPhone UserErr = errors.New("业务操作非法")

	// ErrDriverHaveProcessOrder 当司机还有未支付或未取消订单是返回该错误
	ErrDriverHaveProcessOrder UserErr = errors.New("您有未完成订单")

	// ErrCancelTooMany 当司机当日取消订单次数超过允许值时返回该错误
	ErrCancelTooMany UserErr = errors.New("你今日取消次数过多,不能再发起救援")

	// ErrCancelNoTimes 当司机/技工 当日可取消订单次数超过允许值时返回该错误
	ErrCancelNoTimes UserErr = errors.New("今日已经没有可用取消次数")
)
