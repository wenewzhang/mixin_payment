package models

import (
	"time"
)
// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      gorm.Model
//    }

type Order struct {
	OrderID string `json:"order_id"`
  AssetUUID string `json:"asset_uuid"`
  Amount string `json:"amount"`
  Source string `json:"source"`
}
type OrderTbl struct {
	OrderID string `gorm:"primary_key"`
  AssetUUID string `json:"asset_uuid"`
  Amount string `json:"amount"`
  Source string `json:"source"`
  Status string
  CreatedAt time.Time
  UpdatedAt time.Time
}

type AccountTbl struct {
	OrderID string `gorm:"INDEX"`
	AssetUUID string `json:"asset_uuid"`
  UserID string `gorm:"primary_key"`
  SessionID string
  PinToken string
  PrivateKey string
  Status string
  OpponentID string
  Amount string
  Offset string
  CreatedAt time.Time
  UpdatedAt time.Time
}

type Opponent struct {
  Amount string
  OpponentID string
  TimeStamp string
}
