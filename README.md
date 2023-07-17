# gocrypto 初階交易引擎系統
使用 golang + Gorm(Mysql) + go-redis + grpc + bloom RPC (模擬通信)來實作初階交易引擎。 

使用redis作為開單的緩存來提升開單速度，交易邏輯每個幣種使用兩個queue來儲存使用者買賣單的訂單。

### GOPATH Setup
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

### Compile Module
```
go mod init github.com/mp3ismp3/gocrypto
go mod tidy
```

### RUN
```
go run main.go
```
