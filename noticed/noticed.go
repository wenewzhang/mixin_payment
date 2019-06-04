package main

import (
  // "github.com/wenewzhang/mixin_payment/config"
  mixin "github.com/MooooonStar/mixin-sdk-go/network"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "fmt"
  "time"
)

type AccountTbl struct {
	OrderID string `gorm:"primary_key"`
  UserID string
  SessionID string
  PinToken string
  PrivateKey string
  Status string
  CreatedAt time.Time
  UpdatedAt time.Time
}

// use channel to create wallet,but don't need here!
// c := make(chan mixin.User)
// go createWallet(config.ClientId,config.SessionId, config.PrivateKey, c)
// user := <- c

func readSnapshots( asset string, tm time.Time, userId, sessionId, privateKey string, last chan string) {

  snapData, err := mixin.NetworkSnapshots(asset, tm, true, 30, userId, sessionId, privateKey)
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(snapData)
  last <- "abc"
}

func main() {
  db, err := gorm.Open("sqlite3", "../payment.db")
  if err != nil {
    panic("failed to connect database")
  }
  defer db.Close()

  var accounts  []AccountTbl
  db.Model(&AccountTbl{}).Where("status = ?","pending").Find(&accounts) // find product with id 1
  for _, account := range (accounts) {
    // fmt.Println(account)
    // tm, _:= time.Parse(time.RFC3339Nano,account.CreatedAt.String())
    fmt.Println(account.CreatedAt.Format(time.RFC3339Nano))
    // tm, _:= time.Parse(time.RFC3339Nano,account.CreatedAt.Format(time.RFC3339Nano))
    // fmt.Println(tm)
    c := make(chan string)
    go readSnapshots("", account.CreatedAt, account.UserID,account.SessionID, account.PrivateKey, c)
    last := <- c
    fmt.Println(last)
    // fmt.Println(tm)
  }
}
