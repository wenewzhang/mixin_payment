package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
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

  // Migrate the schema
  // db.AutoMigrate(&Product{})

  // Create
  // db.Create(&Product{Code: "L1212", Price: 1000})

  // Read
  var product OrderTbl
  db.First(&product, 1) // find product with id 1
  fmt.Println(product)
  // db.First(&product, "code = ?", "L1212") // find product with code l1212
  //
  // // Update - update product's price to 2000
  // db.Model(&product).Update("Price", 2000)
  //
  // // Delete - delete product
  // db.Delete(&product)
}
