package model

import (
	"fmt"
	"ocenter/src/model/config"
	"ocenter/src/proto_order"
	gen "ocenter/src/proto_order/rescue"
	"testing"
)

var driverID, oilID int64 = 18598, 18353

var rescueLocation = &proto_order.Location{
	Lat:          99.552061,
	Lng:          199.067990,
	OrderAddress: "中国xxxxxx路",
}

var infos = []*proto_order.OrderInfo{
	&proto_order.OrderInfo{
		GoodsID:  24,
		Name:     "发动机水温高",
		Amount:   2,
		InfoType: 1,
		UnitType: 6,
	},
}

func TestRepairOrderRescue(t *testing.T) {
	checkErr(CloseRescueOrder(driverID))

	d, err := LoadDriver(driverID)
	checkErr(err)

	c, err := config.LoadOtherConf()
	checkErr(err)

	o := NewRepairRescueOrder()
	fsm := NewRepairRescueFSM(o)
	args := &gen.RepairRescueArgs{
		Location: rescueLocation,
	}
	args.OrderInfos = infos

	checkErr(o.TriggerRescueBy(d, fsm, c, args), "发送救援订单失败")
	assert(o.State == 0, fmt.Sprintf("救援订单创建状态错误 = %d, 实际 = %d", 0, o.State))
	id = o.ID
}

func TestRepairOrderRescueDirect(t *testing.T) {
	checkErr(CloseRescueOrder(driverID))

	d, err := LoadDriver(driverID)
	checkErr(err)

	r, err := LoadRepair(oilID)
	checkErr(err)

	c, err := config.LoadOtherConf()
	checkErr(err)

	o := NewRepairRescueOrder()
	fsm := NewRepairRescueFSM(o)
	args := &gen.DirectRescueArgs{
		Distance:   100.0,
		RescueArgs: &gen.RepairRescueArgs{Location: rescueLocation},
	}
	args.GetRescueArgs().OrderInfos = infos

	checkErr(o.TriggerDirectRescueBy(d, r, fsm, c, args), "针对技工发送救援订单失败")
	assert(o.State == 0, fmt.Sprintf("针对技工发送救援订单创建状态错误 = %d, 实际 = %d", 0, o.State))
	id = o.ID
}
