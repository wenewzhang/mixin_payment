package main

import (
  "github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
  "fmt"
  "time"
)

type OrderTbl struct {
	OrderID string `gorm:"primary_key"`
  AssetUUID string `json:"asset_uuid"`
  Amount string `json:"amount"`
  CallBack string `json:"call_back"`
  CreatedAt time.Time
  UpdatedAt time.Time
}

func main() {
  db, err := gorm.Open("sqlite3", "../payment.db")
  if err != nil {
    panic("failed to connect database")
  }
  defer db.Close()
  var users  []OrderTbl

  db.Find(&users) // find product with id 1
  fmt.Println(users)
}
