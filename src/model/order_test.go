package model

import (
	"fmt"
	"log"
)

var id int64

func assert(exp bool, msg string) {
	if !exp {
		log.Fatal(msg)
	}
}

func checkErr(err error, msg ...string) {
	if err != nil {
		panic(fmt.Errorf("%s : %v ", msg, err))
	}
}

// func TestRepairOrderProcess(t *testing.T) {
// 	d, err := LoadDriver(18598)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	r, err := LoadRepair(19581)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	log.Println("老流程维修订单 \n 开始创建...... ")
// 	o, err := placeRepairOrder(d, r)
// 	assert(err == nil, fmt.Sprintf("创建订单错误： %v", err))
// 	assert(o.State == 0, fmt.Sprintf("创建订单状态错误，期望 = %d, 实际 = %d", 0, o.State))
// 	log.Printf("订单创建成功，OrderState = %d, OrderId = %d, Amount = %d", o.State, o.ID, o.Amount)
//
// 	//
// 	log.Println("订单议价开始...... ")
// 	err = bargainRepairOrder(o, r)
// 	assert(err == nil, fmt.Sprintf("订单议价错误： %v", err))
// 	assert(o.State == 1, fmt.Sprintf("订单议价状态错误，期望 = %d, 实际 = %d", 1, o.State))
// 	log.Printf("订单议价成功，OrderState = %d, Amount = %d", o.State, o.Amount)
//
// 	//
// 	log.Println("确认订单开始...... ")
// 	err = confirmedRepairOrder(o, d)
// 	assert(err == nil, fmt.Sprintf("订单确认错误： %v", err))
// 	assert(o.State == 2, fmt.Sprintf("订单确认状态错误，期望 = %d, 实际 = %d", 2, o.State))
// 	log.Printf("订单确认成功，OrderState = %d", o.State)
//
// 	//
// 	log.Println("技工开始出发...... ")
// 	err = runRepairOrder(o, r)
// 	assert(err == nil, fmt.Sprintf("技工出发错误： %v", err))
// 	assert(o.State == 3, fmt.Sprintf("技工出发状态错误，期望 = %d, 实际 = %d", 3, o.State))
// 	log.Printf("技工出发成功，OrderState = %d", o.State)
//
// 	//
// 	log.Println("技工修理中...... ")
// 	err = handRepairOrder(o, r)
// 	assert(err == nil, fmt.Sprintf("技工修理错误： %v", err))
// 	assert(o.State == 4, fmt.Sprintf("技工修理状态错误，期望 = %d, 实际 = %d", 4, o.State))
// 	log.Printf("技工修理成功，OrderState = %d", o.State)
//
// 	//
// 	log.Println("技工完成订单...... ")
// 	err = finishRepairOrder(o, r)
// 	assert(err == nil, fmt.Sprintf("技工完成订单错误： %v", err))
// 	assert(o.State == 5, fmt.Sprintf("技工完成订单状态错误，期望 = %d, 实际 = %d", 5, o.State))
// 	log.Printf("技工完成订单成功，OrderState = %d, Amount = %d", o.State, o.Amount)
//
// 	//
// 	log.Println("订单申请代付...... ")
// 	err = otherPayOrder(o, d)
// 	assert(err == nil, fmt.Sprintf("订单申请代付错误： %v", err))
// 	assert(o.State == 6, fmt.Sprintf("订单申请代付状态错误，期望 = %d, 实际 = %d", 6, o.State))
// 	log.Printf("订单申请代付成功，OrderState = %d", o.State)
//
// 	//
// 	log.Println("订单申请代付...... ")
// 	err = payOrder(o, d)
// 	assert(err == nil, fmt.Sprintf("订单支付错误： %v", err))
// 	assert(o.State == 7, fmt.Sprintf("订单支付状态错误，期望 = %d, 实际 = %d", 7, o.State))
// 	log.Printf("订单支付成功，OrderState = %d", o.State)
//
// 	//
// 	log.Println("订单评价开始...... ")
// 	err = assessOrder(o, d)
// 	assert(err == nil, fmt.Sprintf("订单评价错误： %v", err))
// 	assert(o.State == 8, fmt.Sprintf("订单评价状态错误，期望 = %d, 实际 = %d", 8, o.State))
// 	log.Printf("订单评价成功，OrderState = %d", o.State)
// }
//
// //
// func placeRepairOrder(d *Driver, r *Repair) (*RepairOrder, error) {
// 	params := &proto_repair.RepairCreateArgs{
// 		RepairID: r.ID,
// 		Location: &proto_order.Location{
// 			Lat:          99.552061,
// 			Lng:          199.067990,
// 			OrderAddress: "中国xxxxxx路",
// 		},
// 		Extra:    100,
// 		Distance: 10,
// 	}
// 	info1 := &proto_order.OrderInfo{
// 		GoodsID:  6,
// 		Name:     "发动机",
// 		Amount:   2,
// 		InfoType: 1,
// 		UnitType: 6,
// 		Price:    0,
// 	}
// 	info2 := &proto_order.OrderInfo{
// 		GoodsID:  1,
// 		Name:     "发动机故障报警",
// 		Amount:   1,
// 		InfoType: 2,
// 		UnitType: 6,
// 		ParentID: 6,
// 		Price:    0,
// 	}
// 	info3 := &proto_order.OrderInfo{
// 		GoodsID:  3,
// 		Name:     "检修线路，更换进气压力温度传惑器",
// 		Amount:   1,
// 		InfoType: 3,
// 		UnitType: 6,
// 		ParentID: 1,
// 		Price:    0,
// 	}
// 	infos := []*proto_order.OrderInfo{info1, info2, info3}
// 	params.OrderInfos = infos
// 	o := NewRepairOrder()
// 	fsm := NewRepairFSM(o)
// 	if err := o.TriggerPlaceOrderBy(d, r, fsm, params); err != nil {
// 		return nil, err
// 	}
// 	return o, nil
// }
//
// func bargainRepairOrder(o *RepairOrder, r *Repair) error {
// 	params := &proto_repair.OrderBargainArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxxx路",
// 			},
// 		},
// 		CustomerPrice: &proto_order.CustomerPrice{
// 			Extra:          10,
// 			AccessoriesFee: 11,
// 		},
// 	}
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
// 	infos := []*proto_order.OrderInfo{info1, info2, info3, info4}
// 	params.OrderInfos = infos
// 	fsm := NewRepairFSM(o)
// 	return o.TriggerBargainBy(r, fsm, params)
// }
//
// func confirmedRepairOrder(o *RepairOrder, d *Driver) error {
// 	params := &proto_repair.OrderConfirmeArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxx路",
// 			},
// 		},
// 	}
// 	fsm := NewRepairFSM(o)
// 	return o.TriggerConfirmBy(d, fsm, params)
// }
//
// func runRepairOrder(o *RepairOrder, r *Repair) error {
// 	params := &proto_repair.OrderRunArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxx路",
// 			},
// 		},
// 	}
// 	fsm := NewRepairFSM(o)
// 	return o.TriggerRunBy(r, fsm, params)
// }
//
// func handRepairOrder(o *RepairOrder, r *Repair) error {
// 	params := &proto_repair.OrderHandArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxx华路",
// 			},
// 		},
// 	}
// 	fsm := NewRepairFSM(o)
// 	return o.TriggerHandBy(r, fsm, params)
// }
//
// func finishRepairOrder(o *RepairOrder, r *Repair) error {
// 	params := &proto_repair.OrderFinishArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxxx路",
// 			},
// 		},
// 		CustomerPrice: &proto_order.CustomerPrice{
// 			Extra:          100,
// 			AccessoriesFee: 110,
// 		},
// 	}
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
// 	infos := []*proto_order.OrderInfo{info1, info2, info3, info4}
// 	params.OrderInfos = infos
// 	fsm := NewRepairFSM(o)
// 	return o.TriggerFinishBy(r, fsm, params)
// }
//
// func otherPayOrder(o *RepairOrder, d *Driver) error {
// 	params := &proto_order.OrderOtherPayArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxxx路",
// 			},
// 		},
// 	}
// 	fsm := NewRepairFSM(o)
// 	return o.TriggerOtherPayBy(d, fsm, params)
// }
//
// func payOrder(o *RepairOrder, d *Driver) error {
// 	params := &proto_order.OrderPayArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          119.067990,
// 				OrderAddress: "中国xxxxx路",
// 			},
// 		},
// 	}
// 	fsm := NewRepairFSM(o)
// 	return o.TriggerPayBy(d, fsm, params)
// }
//
// func assessOrder(o *RepairOrder, d *Driver) error {
// 	params := &proto_order.OrderAssessArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxx路",
// 			},
// 		},
// 	}
// 	fsm := NewRepairFSM(o)
// 	return o.TriggerAssessBy(d, fsm, params)
// }

// TestOilOrderProcess 加油订单流程
// func TestOilOrderProcess(t *testing.T) {
//
// 	assert(err == nil, fmt.Sprintf("司机加油订单错误： %v", err))
// 	assert(o.State == 5, fmt.Sprintf("司机加油订单状态错误，期望 = %d, 实际 = %d", 5, o.State))
// 	log.Printf("订单创建成功，OrderState = %d, OrderId = %d, Amount = %d", o.State, o.ID, o.Amount)
// 	//
// 	log.Println("订单申请代付...... ")
// 	err = otherOilPayOrder(o, d)
// 	assert(err == nil, fmt.Sprintf("订单申请代付错误： %v", err))
// 	assert(o.State == 6, fmt.Sprintf("订单申请代付状态错误，期望 = %d, 实际 = %d", 6, o.State))
// 	log.Printf("订单申请代付成功，OrderState = %d", o.State)
//
// 	//
// 	log.Println("订单申请代付...... ")
// 	err = payOilOrder(o, d)
// 	assert(err == nil, fmt.Sprintf("订单支付错误： %v", err))
// 	assert(o.State == 7, fmt.Sprintf("订单支付状态错误，期望 = %d, 实际 = %d", 7, o.State))
// 	log.Printf("订单支付成功，OrderState = %d", o.State)
//
// 	//
// 	log.Println("订单评价开始...... ")
// 	err = assessOilOrder(o, d)
// 	assert(err == nil, fmt.Sprintf("订单评价错误： %v", err))
// 	assert(o.State == 8, fmt.Sprintf("订单评价状态错误，期望 = %d, 实际 = %d", 8, o.State))
// 	log.Printf("加油订单评价成功，OrderState = %d", o.State)
// }

// TestPlaceOilOrderTruckStateErr 加油订单 车辆状态错误测试......
// func TestPlaceOilOrderTruckStateErr(t *testing.T) {
// 	log.Println("测试车辆状态错误")
//
// 	truck, err := LoadCurrentTruck(19602)
// 	checkErr(err, "获取司机当前车辆出错")
// 	assert(truck.State != 2, "车辆状态正常")
// 	assert(truck.checkEffective() != nil, "车辆状态错误测试失败")
// }
//
// // TestPlaceOilOrderTruckVehicleErr 加油订单 车辆类型错误测试......
// func TestPlaceOilOrderTruckVehicleErr(t *testing.T) {
// 	log.Println("测试车辆类型错误")
//
// 	truck, err := LoadCurrentTruck(770)
// 	checkErr(err, "获取司机当前车辆出错")
// 	assert(truck.checkVehicleCarErr() != nil, "车辆类型测试失败")
// }
//

// func createOilOrder() (*OilOrder, error) {
// 	d, err := LoadDriver(18598)
// 	checkErr(err)
//
// 	r, err := LoadRepair(18468)
// 	checkErr(err)
//
// 	params := &proto_oil.OilCreateArgs{
// 		RepairID: r.ID,
// 		Location: &proto_order.Location{
// 			Lat:          99.552061,
// 			Lng:          199.067990,
// 			OrderAddress: "中国xxxx路",
// 			AreaID:       1233456,
// 		},
// 	}
// 	params.OrderInfo = &proto_order.OrderInfo{
// 		GoodsID: 2637,
// 		Name:    "0#",
// 		Amount:  20.0,
// 		Price:   700,
// 	}
//
// 	o := NewOilOrder()
// 	fsm := NewOilFSM(o)
// 	if err := o.TriggerPlaceOrderByDriver(d, r, fsm, params); err != nil {
// 		return nil, err
// 	}
// 	id = o.ID
// 	return o, nil
// }
//
// // TestPlaceOilOrder 加油订单 开始创建......
// func TestPlaceOilOrder(t *testing.T) {
//
// 	log.Println("加油订单 \n 开始创建...... ")
//
// 	o, err := createOilOrder()
// 	checkErr(err, "司机加油订单错误")
//
// 	assert(o.State == 5, fmt.Sprintf("司机加油订单状态错误，期望 = %d, 实际 = %d", 5, o.State))
// 	log.Printf("订单创建成功，OrderState = %d, OrderId = %d, Amount = %d", o.State, o.ID, o.Amount)
// }
//
// //
// func TestOtherOilPayOrder(t *testing.T) {
// 	d, err := LoadDriver(18598)
// 	checkErr(err)
//
// 	o, err := LoadOilOrder(id)
// 	checkErr(err)
// 	fsm := NewOilFSM(o)
// 	params := &proto_order.OrderOtherPayArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxx路",
// 			},
// 		},
// 	}
//
// 	checkErr(o.TriggerOtherPayBy(d, fsm, params), "加油代付处理错误")
// 	assert(o.State == 6, fmt.Sprintf("加油订单代付状态错误 = %d, 实际 = %d", 6, o.State))
// }
//
// //
// func TestPayOilOrder(t *testing.T) {
// 	d, err := LoadDriver(18598)
// 	checkErr(err)
//
// 	o, err := LoadOilOrder(id)
// 	checkErr(err)
// 	fsm := NewOilFSM(o)
// 	params := &proto_order.OrderPayArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxx路",
// 			},
// 		},
// 	}
//
// 	checkErr(o.TriggerPayBy(d, fsm, params), "加油支付处理错误")
// 	assert(o.State == 7, fmt.Sprintf("加油订单支付状态错误 = %d, 实际 = %d", 7, o.State))
// }
//
// //
// func TestAssessOilOrder(t *testing.T) {
// 	d, err := LoadDriver(18598)
// 	checkErr(err)
//
// 	o, err := LoadOilOrder(id)
// 	checkErr(err)
//
// 	fsm := NewOilFSM(o)
// 	params := &proto_order.OrderAssessArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxx路",
// 			},
// 		},
// 	}
//
// 	checkErr(o.TriggerAssessBy(d, fsm, params), "加油评价处理错误")
// 	assert(o.State == 8, fmt.Sprintf("加油订单评价状态错误 = %d, 实际 = %d", 8, o.State))
// }
//
// //
// func TestCloseOilOrder(t *testing.T) {
// 	d, err := LoadDriver(18598)
// 	checkErr(err)
//
// 	_, err = createOilOrder()
// 	checkErr(err, "司机加油订单错误")
//
// 	o, err := LoadOilOrder(id)
// 	checkErr(err)
//
// 	fsm := NewOilFSM(o)
// 	params := &proto_order.OrderCancelArgs{
// 		ProcessBase: &proto_order.ProcessBase{
// 			OrderID: o.ID,
// 			Location: &proto_order.Location{
// 				Lat:          99.552061,
// 				Lng:          199.067990,
// 				OrderAddress: "中国xxxx路",
// 			},
// 		},
// 	}
//
// 	checkErr(o.TriggerDriverClose(d, fsm, params), "加油关闭处理错误")
// 	assert(o.State == 12, fmt.Sprintf("加油订单关闭状态错误 = %d, 实际 = %d", 12, o.State))
// }
