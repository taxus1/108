package model

// var driverID, repairID int64 = 18598, 18353
//
// var rescueLocation = &proto_order.Location{
// 	Lat:          99.552061,
// 	Lng:          199.067990,
// 	OrderAddress: "中国xxxxxx路",
// }
//
// func TestTriggerPlaceOrderBy(t *testing.T) {
// 	d, err := LoadDriver(driverID)
// 	checkErr(err)
//
// 	r, err := LoadRepair(repairID)
// 	checkErr(err)
//
// 	ro, err := LoadOrg(r.OrgID)
// 	checkErr(err)
//
// 	irr, err := config.LoadInsuranceRescueRule()
// 	checkErr(err)
//
// 	o := NewInsuranceRescueOrder()
// 	fsm := NewInsuranceRescueFSM(o)
//
// 	args := &gen.InsuranceRescueArgs{
// 		Distance: 100.0,
// 		Location: rescueLocation,
// 		Insurer:  "中国平安",
// 		Imgs:     []string{"1.jpg", "2.jpg"},
// 	}
//
// 	checkErr(o.TriggerPlaceOrderBy(d, r, fsm, ro, irr, args), "创建保险救援订单失败")
// 	assert(o.State == 5, fmt.Sprintf("保险救援订单创建状态错误 = %d, 实际 = %d", 5, o.State))
// 	id = o.ID
// }

// func TestTriggerDriverCancel(t *testing.T) {
// 	o, err := LoadInsuranceRescueOrder(id)
// 	checkErr(err)
// 	fsm := NewInsuranceRescueFSM(o)
//
// 	d, err := LoadDriver(driverID)
// 	checkErr(err)
//
// 	args := &proto_order.OrderCancelArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          119.067990,
// 				OrderAddress: "中国xxxxx路",
// 			},
// 		},
// 	}
// 	checkErr(o.TriggerDriverCancel(d, fsm, args), "创建保险救援取消失败")
// 	assert(o.State == 9, fmt.Sprintf("保险救援订单取消状态错误 = %d, 实际 = %d", 9, o.State))
// }

// func TestTriggerRejectBy(t *testing.T) {
// 	o, err := LoadInsuranceRescueOrder(id)
// 	checkErr(err)
// 	fsm := NewInsuranceRescueFSM(o)
//
// 	r, err := LoadRepair(repairID)
// 	checkErr(err)
//
// 	args := &proto_order.OrderCancelArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          119.067990,
// 				OrderAddress: "中国xxxxx路",
// 			},
// 		},
// 	}
// 	checkErr(o.TriggerRejectBy(r, fsm, args), "创建保险救援拒绝失败")
// 	assert(o.State == 9, fmt.Sprintf("保险救援订单拒绝状态错误 = %d, 实际 = %d", 9, o.State))
// }
//
// func TestTriggerPayBy(t *testing.T) {
// 	o, err := LoadInsuranceRescueOrder(id)
// 	checkErr(err)
// 	fsm := NewInsuranceRescueFSM(o)
//
// 	r, err := LoadRepair(repairID)
// 	checkErr(err)
//
// 	args := &proto_order.OrderPayArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxx路",
// 			},
// 		},
// 	}
// 	checkErr(o.TriggerPayBy(r, fsm, args), "创建保险救援支付失败")
// 	assert(o.State == 7, fmt.Sprintf("保险救援订单支付状态错误 = %d, 实际 = %d", 7, o.State))
// }
