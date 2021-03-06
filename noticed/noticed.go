package main

import (
  "github.com/wenewzhang/mixin_payment/config"
  "github.com/wenewzhang/mixin_payment/models"
  // "github.com/wenewzhang/mixin_payment/utils"
  mixin "github.com/MooooonStar/mixin-sdk-go/network"
  uuid "github.com/satori/go.uuid"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
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
    log.Println(err)
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
      log.Println("OMG,i find it ----------------------------------------")
      op.OpponentID = v.(map[string]interface{})["opponent_id"].(string)
      op.Amount = v.(map[string]interface{})["amount"].(string)
    }
    // log.Println(v)
    // log.Println(val)
    ctm = v.(map[string]interface{})["created_at"].(string)
  }
  op.TimeStamp = ctm
  log.Println(ctm)
  client <- op
}

func readAssetBalance( asset_uuid string,  userId, sessionId, privateKey string, client chan models.Opponent) {
  var AssetInfo map[string]interface{}
  AssetInfoBytes, err  := mixin.ReadAsset(asset_uuid,
                                         userId,sessionId,privateKey)
  if err != nil { log.Fatal(err) }
  // log.Println(string(AssetInfoBytes))
  if err := json.Unmarshal(AssetInfoBytes, &AssetInfo); err != nil {
      log.Fatal(err)
  }
  log.Println(AssetInfo["data"])
  var op models.Opponent
  if AssetInfo["data"].(map[string]interface{})["balance"].(string) != "0" {
    op.TimeStamp = time.Now().Format(time.RFC3339Nano)
    op.OpponentID = "find"
    op.Amount = AssetInfo["data"].(map[string]interface{})["balance"].(string)
  } else {
    op.TimeStamp = time.Now().Format(time.RFC3339Nano)
    op.OpponentID = ""
    op.Amount = AssetInfo["data"].(map[string]interface{})["balance"].(string)
  }
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
    log.Fatal("failed to connect database")
  }
  defer db.Close()

  var accounts  []models.AccountTbl
  log.Println("run ...")
  for {
    var count int
    db.Model(&models.AccountTbl{}).Where("status = ?","empty").Count(&count)
    for count < 10 {
      user,err := mixin.CreateAppUser("mixin payment", "896400", config.ClientId,
                                     config.SessionId, config.PrivateKey)
      if err != nil {
          log.Fatal("Create account fail, check your config.go!")
      }
      var account models.AccountTbl
      account.OrderID = ""
      account.AssetUUID = ""
      account.UserID = user.UserId
      account.SessionID = user.SessionId
      account.PinToken = user.PinToken
      account.PrivateKey = user.PrivateKey
      account.Status = "empty"
      db.Create(&account)
    }
    //notice
    time.Sleep(config.CheckPendingOrderInterval * time.Second)
    db.Model(&models.AccountTbl{}).Where("status = ?","pending").Find(&accounts) // find product with id 1
    for _, account := range (accounts) {
      log.Println(account)
      // tm, _:= time.Parse(time.RFC3339Nano,account.CreatedAt.String())
      log.Println(account.CreatedAt.Format(time.RFC3339Nano))
      log.Println(time.Since(account.CreatedAt))
      log.Println(account.Offset)
      // m, _ := time.ParseDuration(time.Since(account.CreatedAt))
      // tm, _:= time.Parse(time.RFC3339Nano,account.CreatedAt.Format(time.RFC3339Nano))
      // log.Println(tm)
      c := make(chan models.Opponent)
      // if account.Offset == "" {
      //   go readAssetBalance(account.AssetUUID,  account.UserID,account.SessionID, account.PrivateKey, c)
      // } else {
      go readAssetBalance(account.AssetUUID,  account.UserID,account.SessionID, account.PrivateKey, c)
      // }
      opponent := <- c
      log.Println(opponent.TimeStamp)
      if opponent.OpponentID != "" {
        log.Println("Save the paid infomation---------------" + account.OrderID)
        var order  models.OrderTbl
        if err := db.Model(&models.OrderTbl{}).Where("order_id = ?",account.OrderID).First(&order).Error;err == nil { // find product with id 1
          log.Println(order)
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
          if ( config.MASTER_UUID != "" ) {
             _, err := mixin.Transfer(config.MASTER_UUID, opponent.Amount, order.AssetUUID, order.OrderID,
                                    uuid.Must(uuid.NewV4()).String(),
                                    "896400",account.PinToken,account.UserID,account.SessionID, account.PrivateKey)
             if err != nil {
                     log.Fatal(err)
             }
             log.Println("transfer " + opponent.Amount + " (" + order.AssetUUID + ") to " + config.MASTER_UUID)
          }
        } else
        {
          log.Println("Can not find the order " + account.OrderID + " in table order_tbls!")
        }
      } else {
        if ( (config.OrderExpired * 60) < time.Since(account.CreatedAt).Seconds() ) {
          db.Model(&models.AccountTbl{}).Where("order_id = ?", account.OrderID).Updates(
            map[string]interface{}{"status": "expired"})
        }
      }
      // else {
      //   log.Println("Save timestamp---------------")
      //   db.Model(&models.AccountTbl{}).Where("order_id = ?", account.OrderID).Updates(
      //     map[string]interface{}{"Offset": opponent.TimeStamp})
      // }
      // log.Println(tm)
    }
  } //end of forever for statement
}
