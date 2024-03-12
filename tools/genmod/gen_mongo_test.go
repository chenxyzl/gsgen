package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"gotest/model"
	"math"
	"testing"
)

func TestMongoLoadSave(t *testing.T) {
	var x uint64 = math.MaxUint64 - 1
	var y int64 = int64(x)
	fmt.Println(x)
	fmt.Println(y)

	a := model.TestA{}
	a.SetId(123)
	a.SetA(111)
	a.SetB(222)

	b := model.TestB{}
	b.SetId(456)
	b.SetM("333")
	b.SetN(&a)

	s, e := bson.Marshal(&b)
	if e != nil {
		panic(e)
	}
	b1 := model.TestB{}
	e = bson.Unmarshal(s, &b1)
	if e != nil {
		panic(e)
	}

	//mongo_helper.Connect("mongodb+srv://ichenzhl:Qwert321@cluster0.feqwf3z.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	//defer mongo_helper.Close()
	//col := mongo_helper.GetCol("test", "model")
	//
	//filter := bson.M{"_id": b.GetId()}
	//_, err := col.ReplaceOne(context.TODO(), filter, &b, options.Replace().SetUpsert(true))
	//if err != nil {
	//	log.Fatal(err)
	//}
}
