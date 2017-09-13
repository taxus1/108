package model

import "errors"

// TruckErr 创建订单车辆类错误
type TruckErr error

var (

	// ErrNoTruck 当司机进行维修和保养业务时没有车辆将返回该错误
	ErrNoTruck TruckErr = errors.New("当前无车辆,无法使用维修或保养业务")

	// ErrTruckNotConfirm 车辆没有审核通过错误，当司机进行维修和保养业务时车辆没有审核通过返
	// 回该错误
	ErrTruckNotConfirm TruckErr = errors.New("当前车辆没有审核通过,无法使用维修或保养业务")

	// ErrVehicleTypeIsCar 只有货车才能进行维修或保养业务，当司机车辆为 轿车，面包车 时返回该错误
	ErrVehicleTypeIsCar TruckErr = errors.New("小车,面包车只能使用购买保险业务")

	// ErrVehicleTypeIsBus 只有货车才能进行维修或保养业务，当司机车辆为 客车 时返回该错误
	ErrVehicleTypeIsBus TruckErr = errors.New("客车只能使用加油和保险业务")

	// ErrTruckTransfer 车辆转交中中时创建订单返回该错误
	ErrTruckTransfer TruckErr = errors.New("车辆转交中，暂时无法使用")
)
