#!/bin/sh
# go list -f '{{ join .Deps "\n" }}'|grep github
go get -u github.com/MooooonStar
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/gorilla/mux
go get -u github.com/jinzhu/gorm
go get -u github.com/jinzhu/gorm/dialects/sqlite
go get -u github.com/jinzhu/inflection
go get -u github.com/json-iterator/go
go get -u github.com/mattn/go-sqlite3
go get -u github.com/modern-go/concurrent
go get -u github.com/modern-go/reflect2
go get -u github.com/satori/go.uuid
go get -u github.com/vmihailenco/msgpack
go get -u github.com/vmihailenco/msgpack/codes
go get -u github.com/wenewzhang/mixin_payment/config
go get -u github.com/wenewzhang/mixin_payment/models
go get -u github.com/wenewzhang/mixin_payment/utils
