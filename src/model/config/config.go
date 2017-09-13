package config

import "ocenter/src/db"

// Config 配置
type Config struct {
	ID int64

	//名称
	Name string

	//编码用于程序使用,全局唯一
	Code string

	//内容
	Content string

	//备注
	Remark string

	//消息分类 1短信 2 推送 3 分享 4配置
	Types int32

	//标题
	Title string
}

// LoadConf 加载配置
func loadConf(code string) (*Config, error) {
	sqlStr := `
	SELECT
		id,
		name,
		content,
		remark,
		type,
		title,
		code
	FROM message_config WHERE code = ?
	`
	c := new(Config)
	if err := db.DbSource.QueryRow(sqlStr, code).Scan(c.fileds()); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) fileds() []interface{} {
	return []interface{}{
		&c.ID,
		&c.Name,
		&c.Content,
		&c.Remark,
		&c.Types,
		&c.Title,
		&c.Code,
	}
}
