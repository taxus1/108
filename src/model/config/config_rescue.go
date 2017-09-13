package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"ocenter/src/db"
	"strings"
)

var (
	// ErrNotConfigGrab 当没有配置抢单规则时返回该错误
	ErrNotConfigGrab = errors.New("请先配置推荐技工规则")

	// ErrGrabConfig 当没有配置抢单规则时返回该错误
	ErrGrabConfig = errors.New("抢单规则错误")
)

// RescueBaseConf 救援规则基本配置
type RescueBaseConf struct {
	//技工基础筛选
	MinRepairMiles float32

	MaxRepairMiles float32

	//胎工距离筛选
	MinTireMiles float32

	MaxTireMiles float32
}

func newRescueBaseConf() *RescueBaseConf {
	return &RescueBaseConf{}
}

// RecommendConf 推荐规则
type RecommendConf struct {
	*RescueBaseConf

	//推荐人数
	RecommendNum int32
}

// LoadRecommendConf 加载推荐规则
func LoadRecommendConf() (*RecommendConf, error) {
	r, err := loadRescueConf("RECOMMEND")
	if err != nil {
		return nil, err
	}
	if len(r) == 0 {
		return nil, ErrNotConfigGrab
	}

	c := &RecommendConf{RescueBaseConf: newRescueBaseConf()}
	err = json.Unmarshal(r, c)
	return c, err
}

// GrabConf 抢单规则
type GrabConf struct {
	*RescueBaseConf

	//司机等待时间 分钟
	WaiteTime int32

	//抢单人数
	GrabNum int32
}

// LoadGrabConf 加载推荐规则
func LoadGrabConf() (*GrabConf, error) {
	r, err := loadRescueConf("GRAB")
	if err != nil {
		return nil, ErrGrabConfig
	}
	c := &GrabConf{RescueBaseConf: newRescueBaseConf()}

	if err = json.Unmarshal(r, c); err != nil {
		return nil, ErrGrabConfig
	}
	return c, nil
}

// OtherConf 其它配置
type OtherConf struct {
	//一键救援司机可取消次数 次/天
	DriverCancelTimes int32 `json:"driverCancelTimes"`

	//一键救援技工免费可取消次数 次/天
	RepairCancelTimes int32 `json:"repairCancelTimes"`

	//订单达成后，技工行驶一定公里后司机取消订单需要支付里程费
	NeedPayFeeMiles float32 `json:"needPayFeeMiles"`

	//技工接单后取消订单每次扣除的积分系数
	TakeOffCoefficient float32 `json:"takeOffCoefficient"`

	//技工准备出发-已出发状态变化需要的距离 米
	DepartureMiles int32 `json:"departureMiles"`

	//技工已出发-已到达状态状态变化需要的距离 米
	ArriveMiles int32 `json:"arriveMiles"`
}

// LoadOtherConf 加载推荐规则
func LoadOtherConf() (*OtherConf, error) {
	r, err := loadRescueConf("OTHER")
	if err != nil {
		return nil, err
	}
	c := &OtherConf{}

	if len(r) == 0 {
		return c, nil
	}

	if err = json.Unmarshal(r, c); err != nil {
		return nil, err
	}
	return c, nil
}

// loadRescueConf 加载推荐规则
func loadRescueConf(code string) ([]byte, error) {
	var rule []byte
	sqlStr := "SELECT rule FROM config_dispatcher_orders WHERE code = ?"
	err := db.DbSource.QueryRow(sqlStr, code).Scan(&rule)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return rule, nil
}

// SubsidyConfig 里程补贴
type SubsidyConfig struct {
	Subsidy *Config
	Switch  *Config
}

// LoadSubsidyConf 加载里程补贴配置
func LoadSubsidyConf() (*SubsidyConfig, error) {
	sub, err := loadConf("mileage.config")
	if err != nil {
		return nil, err
	}

	sch, err := loadConf("mileage.switch")
	if err != nil {
		return nil, err
	}

	return &SubsidyConfig{Subsidy: sub, Switch: sch}, nil
}

// Get 获取补贴
func (c *SubsidyConfig) Get() string {
	if c.Switch.Content != "ON" {
		return "0"
	}

	return strings.Split(c.Subsidy.Content, ",")[9]
}
