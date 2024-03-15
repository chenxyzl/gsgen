package test2

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gotest/model"
	"gotest/model/mdata"
	"gotest/tools/genmod/mongo_helper"
	"log"
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

	v := &model.TestA{}
	b.SetC(mdata.NewList[*model.TestA]())
	b.GetC().Append(v)
	//b.GetC().Set(0, v)

	b.SetD(mdata.NewMMap[string, *model.TestA, uint64]())
	b.GetD().Set("1", &model.TestA{})

	b.CleanDirty()

	c := model.TestC{}
	c.SetId(789)
	c.SetX("444")
	c.SetY(&b)
	s, e := bson.Marshal(&c)
	if e != nil {
		panic(e)
	}
	z := model.TestC{}
	e = bson.Unmarshal(s, &z)
	if e != nil {
		panic(e)
	}

	if false {
		mongo_helper.Connect("mongodb+srv://ichenzhl:Qwert321@cluster0.feqwf3z.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
		defer mongo_helper.Close()
		col := mongo_helper.GetCol("test", "model")

		filter := bson.M{"_id": c.GetId()}
		_, err := col.ReplaceOne(context.TODO(), filter, &c, options.Replace().SetUpsert(true))
		if err != nil {
			log.Fatal(err)
		}
		zz := model.TestC{}
		err = col.FindOne(context.TODO(), filter).Decode(&zz)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(zz)
	}
}
