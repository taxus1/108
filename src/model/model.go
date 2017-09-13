package model

// Model 业务模型
type Model struct {
}

var (
	// OrderInfoModel 订单关联的物品业务模型
	OrderInfoModel *orderInfoModel

	// OrderLogModel 订单日志业务模型
	OrderLogModel *orderLogModel

	// OrderModel 订单业务模型
	OrderModel *orderModel

	// TruckModel 车辆业务模型
	TruckModel *truckModel

	// UserModel 用户业务模型
	UserModel *userModel

	// UserSummaryModel 用户统计模型
	UserSummaryModel *userSummaryModel
)

// InitModel 初始化业务模型
func InitModel() {
	OrderInfoModel = &orderInfoModel{}

	OrderLogModel = &orderLogModel{}

	OrderModel = &orderModel{}

	TruckModel = &truckModel{}

	UserModel = &userModel{}

	UserSummaryModel = &userSummaryModel{}
}
