package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"ocenter/src/cst"
	"ocenter/src/db"
	"strings"
)

// truckModel 车辆业务模型
type truckModel Model

// Truck 车辆类
type Truck struct {
	ID int64

	//车辆审核状态
	State int32

	//车辆使用类型
	VehicleType string

	//转交状态(0正常,1转交,2接收)
	Transfer int32

	//车牌号
	PlateNumber string

	//发动机
	Engine string

	//转向
	Trans string

	//驱动形式
	DriverModel string
}

// LoadCurrentTruck 获取司机当前正在使用的车辆
func LoadCurrentTruck(driver int64) (t *Truck, err error) {
	sqlStr := `SELECT t.id, t.status, t.vehicle_type, t.plate_number
          FROM user_truck ut LEFT JOIN truck t ON ut.truck_id = t.id
          WHERE user_id = ? AND is_current = 1
                  AND transfer != 1
                  AND t.deleted = 0
          LIMIT 1
        `
	t = &Truck{}
	err = db.DbSource.QueryRow(sqlStr, driver).Scan(&t.ID, &t.State, &t.VehicleType, &t.PlateNumber)
	if err == sql.ErrNoRows {
		return nil, ErrNoTruck
	}
	return t, nil
}

// GetID 获取ID
func (t *Truck) GetID() int64 {
	if t != nil {
		return t.ID
	}
	return 0
}

// CheckRepaireable 检查车辆使用类型能否使用维修业务
func (t *Truck) CheckRepaireable() error {
	if err := t.checkEffective(); err != nil {
		return err
	}

	if err := t.checkVehicleCarErr(); err != nil {
		return err
	}

	return t.checkVehicleBusErr()
}

// CheckFuelable 检查车辆使用类型能否使用加油业务
func (t *Truck) CheckFuelable() error {
	if err := t.checkEffective(); err != nil {
		return err
	}

	return t.checkVehicleCarErr()
}

// CheckRescueable 检查车辆能否进行救援
func (t *Truck) CheckRescueable() error {
	if err := t.CheckRepaireable(); err != nil {
		return err
	}

	return t.checkTransfer()
}

// TroubleCondition 故障
func (t *Truck) TroubleCondition(buf *bytes.Buffer, troubles *orderInfoTypes) {
	if t == nil {
		return
	}
	for _, v := range *troubles {
		switch v {
		case "6":
			buf.WriteString(fmt.Sprintf(" AND instr(tow.block_engine, '%s') > 0", t.Engine))
		case "7":
			buf.WriteString(fmt.Sprintf(" AND instr(tow.block_trans, '%s') > 0", t.Trans))
		case "8":
			buf.WriteString(fmt.Sprintf(` AND instr(tow.block_brake, '前桥') > 0
						AND instr(tow.block_brand, (SELECT tb.factory_name FROM truck_block_model tbm LEFT JOIN
						truck_brand tb ON tb.brand_id = tbm.brand_id WHERE tbm.model_type = '%s' LIMIT
						1)) > 0`, t.DriverModel))
		case "9":
			buf.WriteString(fmt.Sprintf(` AND instr(tow.block_bridge, '后桥') > 0
				AND instr(tow.block_brand, (SELECT tb.factory_name FROM truck_block_model tbm LEFT JOIN
				truck_brand tb ON tb.brand_id = tbm.brand_id WHERE tbm.model_type = '%s' LIMIT
				1)) > 0`, t.DriverModel))
		case "10":
			buf.WriteString(fmt.Sprintf(` AND instr(tow.block_elec, '电气') > 0
				AND instr(tow.block_brand, (SELECT tb.factory_name FROM truck_block_model tbm LEFT JOIN
				truck_brand tb ON tb.brand_id = tbm.brand_id WHERE tbm.model_type = '%s' LIMIT
				1)) > 0`, t.DriverModel))
		case "11":
			buf.WriteString(fmt.Sprintf(` AND instr(tow.block_body, '车身') > 0
				AND instr(tow.block_brand, (SELECT tb.factory_name FROM truck_block_model tbm LEFT JOIN
				truck_brand tb ON tb.brand_id = tbm.brand_id WHERE tbm.model_type = '%s' LIMIT
				1)) > 0`, t.DriverModel))
		case "12":
			buf.WriteString(fmt.Sprintf(` AND instr(tow.block_maintain, '保养') > 0
				AND instr(tow.block_brand, (SELECT tb.factory_name FROM truck_block_model tbm LEFT JOIN
				truck_brand tb ON tb.brand_id = tbm.brand_id WHERE tbm.model_type = '%s' LIMIT
				1)) > 0`, t.DriverModel))
		}

	}
}

// CheckEffective 检查车辆合法性，有车辆并且审核通过
func (t *Truck) checkEffective() error {
	if t.State != 2 {
		return ErrTruckNotConfirm
	}
	return nil
}

// checkVehicleCarErr 检查车辆类型为轿车，面包车时错误
func (t *Truck) checkVehicleCarErr() error {
	if strings.Contains(t.VehicleType, cst.VechileCar) || strings.Contains(t.VehicleType, cst.VechileSmallBus) {
		return ErrVehicleTypeIsCar
	}

	return nil
}

// checkVehicleBusErr 检查车辆类型为客车时错误
func (t *Truck) checkVehicleBusErr() error {
	if strings.Contains(t.VehicleType, cst.VechileBus) {
		return ErrVehicleTypeIsCar
	}

	return nil
}

func (t *Truck) checkTransfer() error {
	if ok := t.isTransfer(); ok {
		return ErrTruckTransfer
	}
	return nil
}

func (t *Truck) isTransfer() bool {
	return t.Transfer == 1
}
