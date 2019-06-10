package main

import (
    "context"
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
    "fmt"
  	"strings"
	  "encoding/json"
    "encoding/base64"
    "github.com/gorilla/mux"
    "github.com/wenewzhang/mixin_payment/utils"
    "github.com/wenewzhang/mixin_payment/config"
    "github.com/wenewzhang/mixin_payment/models"
    mixin "github.com/MooooonStar/mixin-sdk-go/network"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    // uuid "github.com/satori/go.uuid"
)

// use channel to create wallet,but don't need here!
// c := make(chan mixin.User)
// go createWallet(config.ClientId,config.SessionId, config.PrivateKey, c)
// user := <- c

func createWallet( userId, sessionId, privateKey string, user chan mixin.User) {
  new_user,err := mixin.CreateAppUser("mixin payment", "896400", userId,
                                 sessionId, privateKey)
  if err != nil {
    fmt.Println(err)
  }
  user <- *new_user
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  var order  models.Order
  var orderDB models.OrderTbl
  err := json.NewDecoder(r.Body).Decode(&order) //decode the request body into struct and failed if any error occur
  fmt.Println(order)
  if err != nil {
    utils.Respond(w, utils.Message(false, "Invalid request"))
    return
  }
  db, err := gorm.Open("sqlite3", config.SqlitePath)
  if err != nil {
    panic("failed to connect database")
  }
  defer db.Close()
  fmt.Println(order)
  if config.Assets[order.AssetUUID] == false {
    utils.Respond(w, utils.Message(false, "Unknow Asset UUID!"))
    return
  }
  if db.Model(&models.OrderTbl{}).Where("order_id = ?", order.OrderID).First(&models.OrderTbl{}).RecordNotFound() {
    orderDB.OrderID = order.OrderID
    orderDB.AssetUUID = order.AssetUUID
    orderDB.Amount = order.Amount
    orderDB.Source = order.Source
    db.Create(&orderDB)
    var account models.AccountTbl
    if db.Model(&models.AccountTbl{}).Where("status = ?","empty").First(&account).RecordNotFound() {
      log.Fatal("Please run noticed first!")
      utils.Respond(w, utils.Message(false, "The noticed service won't ready yet!"))
    } else {
      account.OrderID = order.OrderID
      account.AssetUUID = order.AssetUUID
      account.Status = "pending"
      db.Save(&account)
      // payLink := "https://mixin.one/pay?recipient=" +
      //              user.UserId + "&asset=" + order.AssetUUID +
      //              "&amount=" + order.Amount + "&trace=" + uuid.Must(uuid.NewV4()).String() +
      //              "&memo="
      if order.Source == "deposit" {
        UserInfoBytes, err    := mixin.ReadAsset(order.AssetUUID,
                                               account.UserID,account.SessionID,account.PrivateKey)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println(string(UserInfoBytes))
        var UserInfoMap map[string]interface{}
        if err := json.Unmarshal(UserInfoBytes, &UserInfoMap); err != nil {
            panic(err)
        }
        //EOS
        if ( order.AssetUUID == "6cfe566e-4aad-470b-8c9a-2fd35b49c68d" ) {
          log.Println(UserInfoMap["data"].(map[string]interface{})["account_name"])
          log.Println(UserInfoMap["data"].(map[string]interface{})["account_tag"])
          Winfo := utils.EncodeWalletInfo(UserInfoMap["data"].(map[string]interface{})["account_name"].(string),
                                          UserInfoMap["data"].(map[string]interface{})["account_tag"].(string))
          utils.Respond(w, utils.MessagePay(true, "Order has been accepted",Winfo))
        } else {
          log.Println(UserInfoMap["data"].(map[string]interface{})["public_key"])
          Winfo := utils.EncodeWalletInfo(UserInfoMap["data"].(map[string]interface{})["public_key"].(string),
                          "")
          utils.Respond(w, utils.MessagePay(true, "Order has been accepted",Winfo))
        }

      } else {
        payLink := utils.EncodePayurl(account.UserID, order.AssetUUID, order.Amount,order.OrderID)
        fmt.Println(payLink)
        // fmt.Println(user.UserId)
        enUrl := base64.RawURLEncoding.EncodeToString([]byte(payLink))
        utils.Respond(w, utils.MessagePay(true, "Order has been accepted",enUrl))
      }
      return
    }
    } else {
    utils.Respond(w, utils.Message(false, "Order has been denied, because it was existed!"))
    return
  }
}
func checkHandler(w http.ResponseWriter, r *http.Request) {
  var order  models.Order
  err := json.NewDecoder(r.Body).Decode(&order) //decode the request body into struct and failed if any error occur
  fmt.Println(order)
  if err != nil {
    utils.Respond(w, utils.Message(false, "Invalid request"))
    return
  }
  db, err := gorm.Open("sqlite3", config.SqlitePath)
  if err != nil {
    panic("failed to connect database")
  }
  defer db.Close()
  var account  models.AccountTbl
  if db.Model(&models.AccountTbl{}).Where("order_id = ?", order.OrderID).First(&account).RecordNotFound() {
    utils.Respond(w, utils.Message(false, "Order " + order.OrderID + " not existed!"))
  } else {

    // db.Model(&AccountTbl{}).Where("order_id = ?", order.OrderID).First(account)
    fmt.Println(account)
    utils.Respond(w, utils.Message(true, "Order status is " + account.Status))
  }
  return
}
func main() {

    var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()

    // r := mux.NewRouter()
    // Add your routes as needed
    r := mux.NewRouter()
    // r.HandleFunc("/", handler)
    r.HandleFunc("/create_order", createHandler).Methods("POST")
    r.HandleFunc("/check_order", checkHandler).Methods("GET")
    // r.HandleFunc("/articles", handler).Methods("GET")
    // r.HandleFunc("/articles/{id}", handler).Methods("GET", "PUT")
    // r.HandleFunc("/authors", handler).Queries("surname", "{surname}")
    err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
      pathTemplate, err := route.GetPathTemplate()
      if err == nil {
        fmt.Println("ROUTE:", pathTemplate)
      }
      pathRegexp, err := route.GetPathRegexp()
      if err == nil {
        fmt.Println("Path regexp:", pathRegexp)
      }
      queriesTemplates, err := route.GetQueriesTemplates()
      if err == nil {
        fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
      }
      queriesRegexps, err := route.GetQueriesRegexp()
      if err == nil {
        fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
      }
      methods, err := route.GetMethods()
      if err == nil {
        fmt.Println("Methods:", strings.Join(methods, ","))
      }
      fmt.Println()
      return nil
    })

    if err != nil {
      fmt.Println(err)
    }

    http.Handle("/", r)

    srv := &http.Server{
        Addr:         "127.0.0.1:8910",
        // Good practice to set timeouts to avoid Slowloris attacks.
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: r, // Pass our instance of gorilla/mux in.
    }

    // Run our server in a goroutine so that it doesn't block.
    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Println(err)
        }
    }()

//sqlite db
    db, err := gorm.Open("sqlite3", config.SqlitePath)
    if err != nil {
      panic("failed to connect database")
    }
    defer db.Close()

    // Migrate the schema
    db.AutoMigrate(&models.OrderTbl{})
    db.AutoMigrate(&models.AccountTbl{})
    // db.Create(&OrderTbl{})

    c := make(chan os.Signal, 1)
    // We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    // Create a deadline to wait for.
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    // Doesn't block if no connections, but will otherwise wait
    // until the timeout deadline.
    srv.Shutdown(ctx)
    // Optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
    log.Println("shutting down")
    os.Exit(0)
}
