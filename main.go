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
)

type Order struct {
	OrderID string `json:"order_id"`
}


func handler(w http.ResponseWriter, r *http.Request) {
  var order  Order
  err := json.NewDecoder(r.Body).Decode(&order) //decode the request body into struct and failed if any error occur
  fmt.Println(order)
  if err != nil {
    utils.Respond(w, utils.Message(false, "Invalid request"))
    return
  }
  utils.Respond(w, utils.Message(false, "hi,rest api"))
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
    r.HandleFunc("/create_order", handler).Methods("POST")
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
