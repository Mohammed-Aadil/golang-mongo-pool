package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Mohammed-Aadil/foodblog/config/constants"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//MongoPoolError MongoPool error
type MongoPoolError struct {
	When time.Time
	What string
}

func (e *MongoPoolError) Error() string {
	return fmt.Sprintf("at %v, %s",
		e.When, e.What)
}

//IMongoPool mongo db pooling class's interface
type IMongoPool interface {
	Init(string)
	CreateConnection() (*mongo.Client, error)
	GetDatabase() (*mongo.Database, error)
	GetBusyConnectionsCount() int
	GetMaxConnections() int
	GetMinConnections() int
	GetOpenConnections() []*mongo.Client
	GetOpenConnectionsCount() int
	GetPoolName() string
	GetTimeOut() time.Duration
	SetErrorOnBusy()
	SetPoolSize(uint32, uint32)
	TerminateConnection(*mongo.Client)
}

//MongoPool mongo db pooling class
type MongoPool struct {
	name             string
	pool             []*mongo.Client
	minConn          int
	maxConn          int
	raiseErrorOnBusy bool
	poolSize         int
	timeout          time.Duration
}

//Init init MongoPool
func (mp *MongoPool) Init(name string) {
	mp.poolSize = constants.DefaultDBPoolSize
	mp.maxConn = constants.DefaultDBPoolSize
	mp.name = name
	mp.raiseErrorOnBusy = false
	mp.timeout = time.Second * 60 * 5 // 5 minutes
}

func (mp *MongoPool) getContextTimeOut() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), mp.timeout)
	return ctx
}

//CreateConnection create or get a connect
func (mp *MongoPool) CreateConnection() (c *mongo.Client, e error) {
	connections := mp.GetOpenConnections()
	var conn *mongo.Client
	if len(connections) == 0 {
		if mp.raiseErrorOnBusy {
			return nil, &MongoPoolError{
				time.Now(),
				"All connections are busy",
			}
		}
		conn, e = mongo.NewClient(options.Client().ApplyURI(constants.DBUrl))
		if e != nil {
			return nil, e
		}
		e = conn.Connect(mp.getContextTimeOut())
		if e != nil {
			return nil, e
		}
		if mp.poolSize > 0 {
			mp.pool = append(mp.pool, c)
			mp.poolSize--
		}
	} else {
		conn = connections[0]
	}
	return conn, nil
}

//GetDatabase return database obj
func (mp *MongoPool) GetDatabase() (*mongo.Database, error) {
	conn, err := mp.CreateConnection()
	if err != nil {
		return nil, err
	}
	return conn.Database(constants.DBName), err
}

//GetBusyConnectionsCount return no of busy connections
func (mp *MongoPool) GetBusyConnectionsCount() int {
	c := 0
	for _, conn := range mp.pool {
		err := conn.Ping(mp.getContextTimeOut(), readpref.Primary())
		if err != nil {
			c++
		}
	}
	return c
}

//GetMaxConnections return no of max connections
func (mp *MongoPool) GetMaxConnections() int {
	return mp.maxConn
}

//GetMinConnections return no of min connections
func (mp *MongoPool) GetMinConnections() int {
	return mp.minConn
}

//GetOpenConnectionsCount return no of open connections
func (mp *MongoPool) GetOpenConnectionsCount() int {
	c := 0
	for _, conn := range mp.pool {
		err := conn.Ping(mp.getContextTimeOut(), readpref.Primary())
		if err == nil {
			c++
		}
	}
	return c
}

//GetOpenConnections return no of open connections
func (mp *MongoPool) GetOpenConnections() []*mongo.Client {
	var conns []*mongo.Client
	for _, conn := range mp.pool {
		err := conn.Ping(mp.getContextTimeOut(), readpref.Primary())
		if err == nil {
			conns = append(conns, conn)
		}
	}
	return conns
}

//GetPoolName return pool name
func (mp *MongoPool) GetPoolName() string {
	return mp.name
}

//GetTimeOut return timeout
func (mp *MongoPool) GetTimeOut() time.Duration {
	return mp.timeout
}

//SetErrorOnBusy return error when all connection are busy
func (mp *MongoPool) SetErrorOnBusy() {
	mp.raiseErrorOnBusy = true
}

//SetPoolSize set the pool size
func (mp *MongoPool) SetPoolSize(minConn uint32, maxConn uint32) {

}

//TerminateConnection terminate the connection
func (mp *MongoPool) TerminateConnection(conn *mongo.Client) {
	for i, v := range mp.pool {
		if v == conn {
			mp.pool = append(mp.pool[:i], mp.pool[i+1:]...)
			mp.poolSize++
			break
		}
	}
	defer conn.Disconnect(nil)
}

//GetMongoPool get mongo pool instace
func GetMongoPool(name ...string) IMongoPool {
	var mp IMongoPool = new(MongoPool)
	if len(name) == 0 {
		name = append(name, "root")
	}
	mp.Init(name[0])
	return mp
}

//GetCollection return collection obj
func GetCollection(model interface{}) (*mongo.Collection, error) {
	db, err := GetMongoPool().GetDatabase()
	if err != nil {
		return nil, err
	}
	modelString, ok := model.(string)
	if ok {
		return db.Collection(modelString), err
	}
	return db.Collection(reflect.TypeOf(model).Name()), err
}
