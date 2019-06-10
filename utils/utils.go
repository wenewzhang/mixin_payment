package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"github.com/wenewzhang/mixin_payment/models"
	uuid "github.com/satori/go.uuid"
	"encoding/base64"
	"github.com/vmihailenco/msgpack"
)

func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{} {"status" : status, "message" : message}
}

func MessagePay(status bool, message, payurl string) (map[string]interface{}) {
	return map[string]interface{} { "status" : status, "message" : message, "pay": payurl }
}

func Respond(w http.ResponseWriter, data map[string] interface{})  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func EncodePayurl(recipient, asset_id, amount, memo string) (string) {
	baseUrl, _ := url.Parse("https://mixin.one")
	baseUrl.Path += "pay"
	params := url.Values{}
	params.Add("recipient", recipient)
	params.Add("asset", asset_id)
	params.Add("amount", amount)
	params.Add("trace", uuid.Must(uuid.NewV4()).String())
	params.Add("memo", memo)

	// Add Query Parameters to the URL
	baseUrl.RawQuery = params.Encode() // Escape Query Parameters
  return baseUrl.String()
}

func EncodeWalletInfo(public_key, tag string) (string) {
	memoOcean,_ :=
		msgpack.Marshal(models.WalletInfo{
			P: public_key,
			T: tag,
		})
	memoOceanB64 := base64.StdEncoding.EncodeToString(memoOcean)
	return memoOceanB64
}
