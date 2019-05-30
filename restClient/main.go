package main

import (
  "net/http"
  "os"
  "io"
  "encoding/json"
  "bytes"
)
type Order struct {
	OrderID string `json:"order_id"`
}

func main() {
	u := Order{OrderID: "US123"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, _ := http.Post("http://127.0.0.1:8910/create_order", "application/json; charset=utf-8", b)
	io.Copy(os.Stdout, res.Body)
}
