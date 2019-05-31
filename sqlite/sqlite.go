package main

import (
  "github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
  "fmt"
)

type OrderTbl struct {
  gorm.Model
	OrderID string `gorm:"primary_key"`
  AssetUUID string `json:"asset_uuid"`
  Amount string `json:"amount"`
  CallBack string `json:"call_back"`
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
