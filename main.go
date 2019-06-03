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
    "github.com/gorilla/mux"
    "github.com/wenewzhang/mixin_payment/utils"
    "github.com/wenewzhang/mixin_payment/config"
    mixin "github.com/MooooonStar/mixin-sdk-go/network"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Order struct {
	OrderID string `json:"order_id"`
  AssetUUID string `json:"asset_uuid"`
  Amount string `json:"amount"`
  CallBack string `json:"call_back"`
}
type OrderTbl struct {
	OrderID string `gorm:"primary_key"`
  AssetUUID string `json:"asset_uuid"`
  Amount string `json:"amount"`
  CallBack string `json:"call_back"`
  Status string
  CreatedAt time.Time
  UpdatedAt time.Time
}

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

func createWallet( userId, sessionId, privateKey string, user chan mixin.User) {
  new_user,err := mixin.CreateAppUser("mixin payment", "896400", userId,
                                 sessionId, privateKey)
  if err != nil {
    fmt.Println(err)
  }
  user <- *new_user
}

func handler(w http.ResponseWriter, r *http.Request) {
  var order  Order
  var orderDB OrderTbl
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
  fmt.Println(order.OrderID)
  if db.Model(&OrderTbl{}).Where("order_id = ?", order.OrderID).First(&OrderTbl{}).RecordNotFound() {
    orderDB.OrderID = order.OrderID
    orderDB.AssetUUID = order.AssetUUID
    orderDB.Amount = order.Amount
    orderDB.CallBack = order.CallBack
    db.Create(&orderDB)

    // c := make(chan mixin.User)
    // go createWallet(config.ClientId,config.SessionId, config.PrivateKey, c)
    // user := <- c
    user,err := mixin.CreateAppUser("mixin payment", "896400", config.ClientId,
                                   config.SessionId, config.PrivateKey)
    if err != nil {
        utils.Respond(w, utils.Message(false, "Create account fail, check your config.go!"))
    }
    var account AccountTbl
    account.OrderID = order.OrderID
    account.UserID = user.UserId
    account.SessionID = user.SessionId
    account.PinToken = user.PinToken
    account.PrivateKey = user.PrivateKey
    account.Status = "pending"
    db.Create(&account)

    utils.Respond(w, utils.Message(true, "Order has been accepted"))
    return
  } else {
    utils.Respond(w, utils.Message(false, "Order has been denied, because it was existed!"))
    return
  }

}

func main() {

    var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()

    // r := mux.NewRouter()
    // Add your routes as needed
    r := mux.NewRouter()
    // r.HandleFunc("/", handler)
    r.HandleFunc("/create_order", handler).Methods("POST")
    r.HandleFunc("/check_order", handler).Methods("GET")
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
    db.AutoMigrate(&OrderTbl{})
    db.AutoMigrate(&AccountTbl{})
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
