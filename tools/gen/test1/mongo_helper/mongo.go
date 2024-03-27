package mongo_helper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func GetCol(dbname, colName string) *mongo.Collection {
	if mongoClient == nil {
		panic("mongo client not connect")
	}
	return mongoClient.Database(dbname).Collection(colName)
}

func Connect(url string) *mongo.Client {
	// Rest of the code will go here
	// Set client options 设置连接参数
	//clientOptions := options.Client().ApplyURI("mongodb://root:xxxxxxxxxxxxxxxxxxxxxxx@172.20.52.158:41134/?connect=direct;authSource=admin")
	clientOptions := options.Client().ApplyURI(url)

	// Connect to MongoDB 连接数据库
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	// Check the connection 测试连接
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	//
	mongoClient = client
	fmt.Printf("Connected to MongoDB!,url:%v\n", url)
	return mongoClient
}

func Close() {
	if mongoClient != nil {
		mongoClient.Disconnect(context.TODO())
	}
}
