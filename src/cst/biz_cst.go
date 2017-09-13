// Package cst 业务常量
package cst

// TruckCheckPass 车辆审核通过状态
const TruckCheckPass = 2

// VehicleTypes 行驶证上车辆使用类型
const (
	VechileCar = "轿车"

	VechileSmallBus = "面包车"

	VechileBus = "客车"
)

// 技工状态
const (
	//休息
	RepairRest = 4

	//工作中
	RepairWork = 5

	//繁忙
	RepairBusy = 7
)

// 订单类型 1 油 2 轮胎 3 买新胎 4 补旧胎 5集团加油
// 6发动机 7离合/变速 8转向制动 9车桥悬挂 10电气 11车身 12保养
// 13自营油 14加油站油 15加盟油 17维修 18商城机油订单  app不使用
// 19 保险救援 20 保证金
const (
	OrderTypeAll             = 0
	OrderTypeOil             = 1
	OrderTypeTire            = 2
	OrderTypeTireBuy         = 3
	OrderTypeTireRepair      = 4
	OrderTypeOilGroup        = 5
	OrderTypeEngine          = 6
	OrderTypeTrans           = 7
	OrderTypeDrive           = 8
	OrderTypeBridge          = 9
	OrderTypeElec            = 10
	OrderTypeBody            = 11
	OrderTypeMaintain        = 12
	OrderTypeOilSelf         = 13
	OrderTypeOilStation      = 14
	OrderTypeOilOutside      = 15
	OrderTypeTruckRepair     = 17
	OrderTypeMallMachineOil  = 18
	OrderTypeInsuranceRescue = 19
	OrderTypeOrgBond         = 20
)
