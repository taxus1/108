package model

import (
	"database/sql"
	"ocenter/src/db"
	"strings"
	"time"
)

var imgCols = []string{
	"order_id",
	"img",
	"create_time",
	"update_time",
	"status",
}

// OrderImage 订单车头照
type OrderImage struct {
	ID int64

	OrderID int64

	Img string

	State int32

	CreateTime time.Time

	UpdateTime time.Time
}

// NewOrderImage 新建订单车头照
func NewOrderImage(orderID int64) *OrderImage {
	return &OrderImage{
		ID:         orderID,
		State:      0,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

// LoadOrderImage 获取订单车头照
func LoadOrderImage(orderID int64) (*OrderImage, error) {
	sqlStr := `
    SELECT
      id,
      img,
      create_time,
      update_time,
      status
    FROM order_images WHERE order_id = ?
  	`
	oi := new(OrderImage)
	err := db.DbSource.QueryRow(sqlStr, orderID).Scan(
		&oi.ID,
		&oi.Img,
		&oi.CreateTime,
		&oi.UpdateTime,
		&oi.State,
	)
	if err != nil {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return NewOrderImage(orderID), nil
	}
	return oi, nil
}

// SetImg 设置照片
func (oi *OrderImage) SetImg(imgs []string) {
	oi.Img = strings.Join(imgs, ",")
}

// GetImg 获取照片
func (oi *OrderImage) GetImg() []string {
	return strings.Split(oi.Img, ",")
}

// Upload 上传图片
func (oi *OrderImage) Upload(tx *sql.Tx) error {
	if oi.ID != 0 {
		return oi.update(tx)
	}
	return oi.save(tx)
}

func (oi *OrderImage) update(tx *sql.Tx) error {
	sqlStr := "UPDATE order_images SET img = ? WHERE id = ?"
	if _, err := db.UpdateTx(tx, sqlStr, oi.Img, oi.ID); err != nil {
		return err
	}
	return nil
}

func (oi *OrderImage) save(tx *sql.Tx) error {
	sqlStr := db.BuildInsert("order_images", imgCols, 1)
	lastID, err := db.SaveTx(tx, sqlStr, oi.values())
	if err != nil {
		return err
	}

	oi.ID = lastID
	return nil
}

func (oi *OrderImage) values() []interface{} {
	return []interface{}{
		oi.OrderID,
		oi.Img,
		oi.State,
		oi.CreateTime,
		oi.UpdateTime,
	}
}
