package test

import (
	"context"
	"fmt"
	"gen_tools/model"
	"gen_tools/model/mdata"
	"gen_tools/test/mongo_helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	b.SetC(mdata.NewList[*model.TestA]())
	v1 := &model.TestA{}
	v1.SetId(1)
	v1.SetA(2)
	v1.SetB(3)
	b.GetC().Append(v1)
	v2 := &model.TestA{}
	v2.SetId(11)
	v2.SetA(12)
	v2.SetB(13)
	b.GetC().Append(v2)
	//b.GetC().Set(0, v)

	b.SetD(mdata.NewMMap[string, *model.TestA]())
	v11 := &model.TestA{}
	v11.SetId(100)
	v11.SetA(101)
	v11.SetB(102)
	b.GetD().Set("100", v11)
	v22 := &model.TestA{}
	v22.SetId(110)
	v22.SetA(111)
	v22.SetB(112)
	b.GetD().Set("110", v22)
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
		mongo_helper.Connect("") //todo 换成自己的mongo地址测试
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
