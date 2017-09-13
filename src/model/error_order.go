package model

import "errors"

// OrderErr 订单相关错误
type OrderErr error

var (

	// ErrHasProcessOrder 当司机或技工有未完成订单时返回该错误
	// 有进行中的订单就是有订单状态为未完成（!5）的订单
	ErrHasProcessOrder OrderErr = errors.New("上一订单尚未完成")

	// ErrOrderAmount 当订单金额小于0时返回错误
	ErrOrderAmount OrderErr = errors.New("订单总金额错误，请检查优惠金额和辅料费")

	// ErrOrderPaid 当重复支付订单时返回该错误
	ErrOrderPaid OrderErr = errors.New("订单已经支付")

	// ErrOrderInvalid 当订单已经失效时还要进行订单操作 返回该错误
	ErrOrderInvalid OrderErr = errors.New("订单已关闭")

	// ErrOrderCanceld 订单取消时还要进行订单操作 返回该错误
	ErrOrderCanceld OrderErr = errors.New("订单已取消")

	// ErrRescueGrabed 救援订单被别人抢走时返回该错误
	ErrRescueGrabed OrderErr = errors.New("您下手太慢了, 该订单已被抢走,下次要尽快下手哟~")

	// ErrRescueCanceled 救援订单取消时返回该错误
	ErrRescueCanceled OrderErr = errors.New("司机取消了该订单,非常抱歉")

	// ErrRescueExpired 救援订单抢单超时返回该错误
	ErrRescueExpired OrderErr = errors.New("该订单已过期,请及时接单哟~")

	// ErrOrgHasInsuranceOrder 维修厂有未处理保险救援订单时返回该错误
	ErrOrgHasInsuranceOrder OrderErr = errors.New("维修厂正在处理其他救援,稍后再试")

	// ErrOrderInPay 订单支付中时取消订单返回给错误
	ErrOrderInPay OrderErr = errors.New("维修厂正在处理您的救援订单,不能取消订单")
)
