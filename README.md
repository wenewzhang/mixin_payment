# mixin_payment
Usage:
```bash
//run server
go run main.go

//run client
cd restClient
go run main.go
```

api:
- check_order
```bash
curl -X POST -i 'http://127.0.0.1:8910/check_order' --data '{"order_id":"1211321"}'
```

- create_order
```bash
curl -X POST -i 'http://127.0.0.1:8910/create_order' --data '{"order_id":"1211322","asset_uuid":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","amount":"1","call_back":""}'

curl -X POST -i 'http://127.0.0.1:8910/create_order' --data '{"order_id":"1211323","asset_uuid":"965e5c6e-434c-3fa9-b780-c50f43cd955c","amount":"1","call_back":""}'
```
