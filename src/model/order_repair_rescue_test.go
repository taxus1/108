package model

// var driverID, oilID int64 = 18598, 18353
//
// var rescueLocation = &proto_order.Location{
// 	Lat:          99.552061,
// 	Lng:          199.067990,
// 	OrderAddress: "中国xxxxxx路",
// }
//
// var infos = []*proto_order.OrderInfo{
// 	&proto_order.OrderInfo{
// 		GoodsID:  24,
// 		Name:     "发动机水温高",
// 		Amount:   2,
// 		InfoType: 1,
// 		UnitType: 6,
// 	},
// }

// func TestRepairOrderRescue(t *testing.T) {
// 	checkErr(CloseRescueOrder(driverID))
//
// 	d, err := LoadDriver(driverID)
// 	checkErr(err)
//
// 	c, err := config.LoadOtherConf()
// 	checkErr(err)
//
// 	o := NewRepairRescueOrder()
// 	fsm := NewRepairRescueFSM(o)
// 	args := &gen.RepairRescueArgs{
// 		Location: rescueLocation,
// 	}
// 	args.OrderInfos = infos
//
// 	checkErr(o.TriggerRescueBy(d, fsm, c, args), "发送救援订单失败")
// 	assert(o.State == 0, fmt.Sprintf("救援订单创建状态错误 = %d, 实际 = %d", 0, o.State))
// 	id = o.ID
// 	fmt.Println(id)
// }

// func TestTriggerGrabBy(t *testing.T) {
// 	o, err := LoadRepairRescueOrder(id)
// 	checkErr(err)
//
// 	g, err := LoadRescue(o.ID)
// 	checkErr(err)
//
// 	c, err := config.LoadGrabConf()
// 	checkErr(err)
//
// 	r, err := LoadRepair(oilID)
// 	checkErr(err)
//
// 	fsm := NewRepairRescueFSM(o)
// 	args := &gen.RescueGrabArgs{
// 		Distance: 100.0,
// 		Location: rescueLocation,
// 	}
// 	checkErr(o.TriggerGrabBy(r, fsm, c, g, args), "技工抢单失败")
// 	assert(o.State == 3, fmt.Sprintf("技工抢单状态错误 = %d, 实际 = %d", 3, o.State))
// }

// func TestTriggerCompleteBy(t *testing.T) {
// 	o, err := LoadRepairRescueOrder(1000295)
// 	checkErr(err)
//
// 	r, err := LoadRepair(oilID)
// 	checkErr(err)
//
// 	fsm := NewRepairRescueFSM(o)
// 	args := &gen.RescueCompleteArgs{
// 		Distance: 100.0,
// 		Image:    []string{"1.jpg", "2.jpg"},
// 		CustomerPrice: &proto_order.CustomerPrice{
// 			Extra:          10,
// 			AccessoriesFee: 11,
// 		},
// 		Location: rescueLocation,
// 	}
//
// 	info1 := &proto_order.OrderInfo{
// 		GoodsID:  1,
// 		Name:     "保养",
// 		Amount:   0,
// 		InfoType: 1,
// 		UnitType: 12,
// 		Price:    0,
// 	}
// 	info2 := &proto_order.OrderInfo{
// 		GoodsID:  -1,
// 		Name:     "发动机保养",
// 		Amount:   0,
// 		InfoType: 2,
// 		UnitType: 0,
// 		ParentID: 1,
// 		Price:    0,
// 	}
// 	info3 := &proto_order.OrderInfo{
// 		GoodsID:  496,
// 		Name:     "更换机油",
// 		Amount:   1,
// 		InfoType: 3,
// 		UnitType: 0,
// 		ParentID: -1,
// 		Price:    0,
// 	}
// 	info4 := &proto_order.OrderInfo{
// 		GoodsID:  582,
// 		Name:     "机油（大）",
// 		Amount:   1,
// 		InfoType: 4,
// 		UnitType: 0,
// 		ParentID: 496,
// 		Price:    680,
// 	}
// 	args.OrderInfos = []*proto_order.OrderInfo{info1, info2, info3, info4}
// 	checkErr(o.TriggerCompleteBy(r, fsm, args), "技工完善订单失败")
// 	assert(o.State == 5, fmt.Sprintf("技工完善订单状态错误 = %d, 实际 = %d", 5, o.State))
// }

// func TestRepairOrderRescueDirect(t *testing.T) {
// 	checkErr(CloseRescueOrder(driverID))
//
// 	d, err := LoadDriver(driverID)
// 	checkErr(err)
//
// 	r, err := LoadRepair(oilID)
// 	checkErr(err)
//
// 	c, err := config.LoadOtherConf()
// 	checkErr(err)
//
// 	o := NewRepairRescueOrder()
// 	fsm := NewRepairRescueFSM(o)
// 	args := &gen.DirectRescueArgs{
// 		Distance:   100.0,
// 		RescueArgs: &gen.RepairRescueArgs{Location: rescueLocation},
// 	}
// 	args.GetRescueArgs().OrderInfos = infos
//
// 	checkErr(o.TriggerDirectRescueBy(d, r, fsm, c, args), "针对技工发送救援订单失败")
// 	assert(o.State == 0, fmt.Sprintf("针对技工发送救援订单创建状态错误 = %d, 实际 = %d", 0, o.State))
// 	id = o.ID
// }

// func TestTriggerAcceptBy(t *testing.T) {
// 	o, err := LoadRepairRescueOrder(id)
// 	checkErr(err)
// 	fsm := NewRepairRescueFSM(o)
//
// 	r, err := LoadRepair(oilID)
// 	checkErr(err)
//
// 	args := &gen.RescueAcceptArgs{
// 		AcceptArgs: &gen.RescueGrabArgs{
// 			Distance: 100.0,
// 			Location: rescueLocation,
// 		},
// 	}
// 	checkErr(o.TriggerAcceptBy(r, fsm, args), "技工接受救援订单失败")
// 	assert(o.State == 3, fmt.Sprintf("技工接受救援订单状态错误 = %d, 实际 = %d", 3, o.State))
// }
//
// func TestTriggerDriverCancel(t *testing.T) {
// 	o, err := LoadRepairRescueOrder(id)
// 	checkErr(err)
// 	fsm := NewRepairRescueFSM(o)
//
// 	d, err := LoadDriver(driverID)
// 	checkErr(err)
//
// 	r, err := LoadRepair(oilID)
// 	checkErr(err)
//
// 	c, err := config.LoadOtherConf()
// 	checkErr(err)
//
// 	res, err := LoadRescue(o.ID)
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
// 	checkErr(o.TriggerDriverCancel(d, r, fsm, c, res, args), "司机取消救援订单失败")
// 	assert(o.State == 9, fmt.Sprintf("机取消救援订单状态错误 = %d, 实际 = %d", 9, o.State))
// }

// func TestTriggerRepairCancel(t *testing.T) {
// 	o, err := LoadRepairRescueOrder(id)
// 	checkErr(err)
// 	fsm := NewRepairRescueFSM(o)
//
// 	r, err := LoadRepair(oilID)
// 	checkErr(err)
//
// 	c, err := config.LoadOtherConf()
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
// 	checkErr(o.TriggerRepairCancel(r, fsm, c, args), "技工取消救援订单失败")
// 	assert(o.State == 9, fmt.Sprintf("技工消救援订单状态错误 = %d, 实际 = %d", 9, o.State))
// }

// func TestTriggerRejectBy(t *testing.T) {
// 	o, err := LoadRepairRescueOrder(id)
// 	checkErr(err)
// 	fsm := NewRepairRescueFSM(o)
//
// 	r, err := LoadRepair(oilID)
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
// 	checkErr(o.TriggerRejectBy(r, fsm, args), "技工拒绝救援订单失败")
// 	assert(o.State == 9, fmt.Sprintf("技工拒绝援订单状态错误 = %d, 实际 = %d", 9, o.State))
// }
