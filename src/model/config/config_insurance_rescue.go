package config

import (
	"encoding/json"
	"time"
)

// InsuranceRescueRule 保险救援活动规则
type InsuranceRescueRule struct {
	*Config
	Activity bool  `json:"isActivity"`
	Type     int32 `json:"type"`
	Number   int32 `json:"number"`
	Amount   int64 `json:"amount"`
}

// LoadInsuranceRescueRule 加载保险救援活动配置
func LoadInsuranceRescueRule() (*InsuranceRescueRule, error) {
	r, err := loadConf("rescue.servicecharge")
	if err != nil {
		return nil, err
	}

	c := &InsuranceRescueRule{Config: r}
	err = json.Unmarshal([]byte(r.Remark), c)
	return c, err
}

// HasActivity 是否有活动
func (i *InsuranceRescueRule) HasActivity() bool {
	return i.Activity
}

// IsLimitTimeActivity 是否限时活动
func (i *InsuranceRescueRule) IsLimitTimeActivity() bool {
	return i.Type == 0
}

// Expire 按时间限制活动时 是否过期
func (i *InsuranceRescueRule) Expire(t time.Time) bool {
	return time.Now().Unix() > t.Unix()+int64(i.Number*24*60*60)
}

// Overrun 按次数限制活动时 是否超限
func (i *InsuranceRescueRule) Overrun(n int32) bool {
	return i.Number < n
}
