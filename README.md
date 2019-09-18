# golang-mongo-pool

This repository is uploaded for academic learning of database pooling class. I have used official golang mongo-driver for mongodb.
For installation please go to your GOPATH and enter the following command:

```shell
go get github.com/Mohammed-Aadil/golang-mongo-pool
```

For pooling following functions are required in class:

1. `Init(string)`
2. `CreateConnection() (*mongo.Client, error)`
3. `GetDatabase() (*mongo.Database, error)`
4. `GetBusyConnectionsCount() int`
5. `GetMaxConnections() int`
6. `GetMinConnections() int`
7. `GetOpenConnections() []*mongo.Client`
8. `GetOpenConnectionsCount() int`
9. `GetPoolName() string`
10. `GetTimeOut() time.Duration`
11. `SetErrorOnBusy()`
12. `SetPoolSize(uint32, uint32)`
13. `TerminateConnection(*mongo.Client)`

## Init(string)

This function is use for initializing important config of mongo db eg: connection client

## CreateConnection() (*mongo.Client, error)

This function is used as getObject of pooling class. It will return new connection if required on calling.

## GetDatabase() (*mongo.Database, error)

This function is call to initialize database and return db instance from mongo client.

## GetBusyConnectionsCount() int

This Function return no of busy connection count.

## GetMaxConnections() int

This function return count of maximum allowed connections in pool.

## GetMinConnections() int

This function return count of minimum allowed connections in pool.

## GetOpenConnections() []*mongo.Client

This function return all the open connection instances in pool.

## GetOpenConnectionsCount() int

This function return count of open connections in pool.

## GetPoolName() string

This function return pool name. Default is `root`

## GetTimeOut() time.Duration

This function return timeout duration of connection in pool.

## SetErrorOnBusy()

If this function is called then if all connections are busy `CreateConnection()` function will return `MongoPoolError`.

## SetPoolSize(uint32, uint32)

This will resize mongo pool without affecting currently open connection. **This function still need to be implemented. All PR for this are welcome.**

## TerminateConnection(*mongo.Client)

This will terminate the connection and release same object from pool.
