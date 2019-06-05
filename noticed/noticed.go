package main

import (
  "github.com/wenewzhang/mixin_payment/config"
  "github.com/wenewzhang/mixin_payment/models"
  mixin "github.com/MooooonStar/mixin-sdk-go/network"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "fmt"
  "time"
  "encoding/json"
  "log"
  "strconv"
)

// use channel to create wallet,but don't need here!
// c := make(chan mixin.User)
// go createWallet(config.ClientId,config.SessionId, config.PrivateKey, c)
// user := <- c

func readSnapshots( asset string, tm time.Time, userId, sessionId, privateKey string, client chan models.Opponent) {

  snapData, err := mixin.NetworkSnapshots(asset, tm, true, 100, userId, sessionId, privateKey)
  if err != nil {
    fmt.Println(err)
  }
  var snapInfo map[string]interface{}
  err = json.Unmarshal([]byte(snapData), &snapInfo)
  if err != nil {
      log.Fatal(err)
  }
  var ctm string
  var op models.Opponent
  for _, v := range (snapInfo["data"].([]interface{})) {
    if v.(map[string]interface{})["opponent_id"] != nil {
      fmt.Println("OMG,i find it ----------------------------------------")
      op.OpponentID = v.(map[string]interface{})["opponent_id"].(string)
      op.Amount = v.(map[string]interface{})["amount"].(string)
    }
    // fmt.Println(v)
    // fmt.Println(val)
    ctm = v.(map[string]interface{})["created_at"].(string)
  }
  op.TimeStamp = ctm
  fmt.Println(ctm)
  client <- op
}

/*---my snapshots format----*/
// {"amount"=>"0.00013147", "asset"=>{"asset_id"=>"c6d0c728-2624-429b-8e0d-d9d19b6592fa",
// "asset_key"=>"c6d0c728-2624-429b-8e0d-d9d19b6592fa",
//  "chain_id"=>"c6d0c728-2624-429b-8e0d-d9d19b6592fa",
//  "icon_url"=>"https://images.mixin.one/HvYGJsV5TGeZ-X9Ek3FEQohQZ3fE9LBEBGcOcn4c4BNHovP4fW4YB97Dg5LcXoQ1hUjMEgjbl1DPlKg1TW7kK6XP=s128",
//  "name"=>"Bitcoin", "symbol"=>"BTC", "type"=>"asset"},
//  "created_at"=>"2019-05-23T09:48:04.582099Z",
//  "data"=>"hqFDzQPooVCnNzU2Mi45MaFGqzAuMDAwMDAwMjY0okZBxBDG0McoJiRCm44N2dGbZZL6oVShUqFPxBDEACoH8bFDObzOJcNDiF5S",
//   "opponent_id"=>"61103d28-3ac2-44a2-ae34-bd956070dab1",
//   "snapshot_id"=>"dabcad80-4352-4d24-8599-73d374dfaebd", "source"=>"TRANSFER_INITIALIZED",
//   "trace_id"=>"bdc79adc-f2b3-4eeb-953d-01b476f91322",
//  "type"=>"snapshot",
//  "user_id"=>"5e4ad097-21e8-3f6b-98f7-9dc74dd99f77"}

func main() {
  db, err := gorm.Open("sqlite3", "../payment.db")
  if err != nil {
    panic("failed to connect database")
  }
  defer db.Close()

  var accounts  []models.AccountTbl
  fmt.Println("run ...")
  for {
    time.Sleep(config.CheckPendingOrderInterval * time.Second)
    db.Model(&models.AccountTbl{}).Where("status = ?","pending").Find(&accounts) // find product with id 1
    for _, account := range (accounts) {
      // fmt.Println(account)
      // tm, _:= time.Parse(time.RFC3339Nano,account.CreatedAt.String())
      fmt.Println(account.CreatedAt.Format(time.RFC3339Nano))
      fmt.Println(time.Since(account.CreatedAt))
      fmt.Println(account.Offset)
      if ( (config.OrderExpired * 60) < time.Since(account.CreatedAt) ) {
        db.Model(&models.AccountTbl{}).Where("order_id = ?", account.OrderID).Updates(
          map[string]interface{}{"status": "expired"})
        continue
      }
      // tm, _:= time.Parse(time.RFC3339Nano,account.CreatedAt.Format(time.RFC3339Nano))
      // fmt.Println(tm)
      c := make(chan models.Opponent)
      if account.Offset == "" {
        go readSnapshots("", account.CreatedAt, account.UserID,account.SessionID, account.PrivateKey, c)
      } else {
        tmOffset, _ := time.Parse(time.RFC3339Nano,account.Offset)
        go readSnapshots("", tmOffset, account.UserID,account.SessionID, account.PrivateKey, c)
      }
      opponent := <- c
      fmt.Println(opponent.TimeStamp)
      if opponent.OpponentID != "" {
        fmt.Println("Save the paid infomation---------------" + account.OrderID)
        var order  models.OrderTbl
        if err := db.Model(&models.OrderTbl{}).Where("order_id = ?",account.OrderID).First(&order).Error;err == nil { // find product with id 1
          fmt.Println(order)
          rAmount, _ := strconv.ParseFloat(order.Amount, 64)
          tAmount, _ := strconv.ParseFloat(opponent.Amount, 64)
          if rAmount <= tAmount {
            db.Model(&models.AccountTbl{}).Where("order_id = ?", account.OrderID).Updates(
              map[string]interface{}{"opponent_id": opponent.OpponentID,
                                    "amount":opponent.Amount,
                                    "status":"paid"})
          } else {
            db.Model(&models.AccountTbl{}).Where("order_id = ?", account.OrderID).Updates(
              map[string]interface{}{"opponent_id": opponent.OpponentID,
                                     "amount":opponent.Amount,
                                     "status":"partial-paid"})
          }
        } else
        {
          fmt.Println("Can not find the order " + account.OrderID + " in table order_tbls!")
        }
      } else {
        fmt.Println("Save timestamp---------------")
        db.Model(&models.AccountTbl{}).Where("order_id = ?", account.OrderID).Updates(
          map[string]interface{}{"Offset": opponent.TimeStamp})
      }
      // fmt.Println(tm)
    }
  } //end of forever for statement
}
