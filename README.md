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
curl -X GET -i 'http://127.0.0.1:8910/check_order' --data '{"order_id":"1211323"}'
```

- create_order
```bash
curl -X POST -i 'http://127.0.0.1:8910/create_order' --data '{"order_id":"1211322","asset_uuid":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","amount":"1","call_back":""}'


curl -X POST -i 'http://127.0.0.1:8910/create_order' --data '{"order_id":"1211330","asset_uuid":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","amount":"0.006","source":"mixin"}'

curl -X POST -i 'http://127.0.0.1:8910/create_order' --data '{"order_id":"1211329","asset_uuid":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","amount":"0.006","source":"deposit"}'

```

install in ubuntu 16.04 LTS:
```bash
sudo apt-get update
sudo apt-get -y upgrade
sudo apt-get install g++ make
wget https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz
sudo tar -xvf go1.12.5.linux-amd64.tar.gz
sudo mv go /usr/local

vi ~/.bashrc
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
source ~/.bashrc
go get github.com/wenewzhang/mixin_payment

cd noticed
make
./noticed &

cd ../
make
./mixin_payment &
```
