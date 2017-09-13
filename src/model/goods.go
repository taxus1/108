package model

import (
	"database/sql"
	"ocenter/src/db"
	"time"
)

var goodsCols = []string{
	"id",
	"name",
	"brand_id",
	"model_id",
	"price",
	"info",
	"org_id",
	"is_goods",
	"unit",
	"img",
	"create_time",
	"stripe_id",
	"status",
	"update_time",
	"benefit",
}

// GoodsSchema 物品/商品
type GoodsSchema struct {
	ID         int64
	Name       string
	BrandID    int64
	ModelID    int64
	Price      int64
	Info       string
	OrgID      int64
	IsGoods    bool
	Unit       string
	Img        string
	CreateTime time.Time
	StripeID   int64
	Status     int32
	UpdateTime time.Time
	Benefit    int64
}

// Goods 物品/商品
type Goods struct {
	*GoodsSchema
}

// LoadGoods 获取物品
func LoadGoods(orgID int64) (*Goods, error) {
	sqlStr := `
        SELECT t.*
        FROM goods t
        LEFT JOIN goods_brand gb ON t.brand_id = gb.id
        WHERE t.org_id = ? AND gb.code = 'TIRE_MILE' AND  status = 1
        LIMIT 1
  `
	g := &Goods{GoodsSchema: &GoodsSchema{}}
	err := db.DbSource.QueryRow(sqlStr, orgID).Scan(g.fields()...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return g, nil
}

func (g *Goods) fields() []interface{} {
	return []interface{}{
		&g.ID,
		&g.Name,
		&g.BrandID,
		&g.ModelID,
		&g.Price,
		&g.Info,
		&g.OrgID,
		&g.IsGoods,
		&g.Unit,
		&g.Img,
		&g.CreateTime,
		&g.StripeID,
		&g.Status,
		&g.UpdateTime,
		&g.Benefit,
	}
}
