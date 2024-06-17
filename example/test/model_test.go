package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chenxyzl/gsgen/example/nest"
	"github.com/chenxyzl/gsgen/gsmodel"
	"github.com/chenxyzl/gsgen/test/mongo_helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

var mongoUrl = "" //todo 换成自己的mongo地址测试

func getTestC(id uint64) *nest.TestC {
	c := &nest.TestC{}
	c.SetId(id)
	c.SetA("c_a")
	c.SetB(&nest.TestB{})
	c.GetB().SetId(2)
	c.GetB().SetAa("c_b_2")
	c.GetB().SetBb(&nest.TestA{})
	c.GetB().GetBb().SetId(100)
	c.GetB().GetBb().SetAaa(101)
	c.GetB().GetBb().SetBbb(102)
	c.GetB().GetBb().SetCcc("103")
	c.GetB().SetCc(gsmodel.NewDList[*nest.TestA]())
	c.GetB().GetCc().Append(&nest.TestA{})
	c.GetB().GetCc().Get(0).SetId(3)
	c.GetB().GetCc().Get(0).SetAaa(4)
	c.GetB().GetCc().Get(0).SetBbb(5)
	c.GetB().GetCc().Get(0).SetCcc("6")
	c.GetB().GetCc().Append(&nest.TestA{})
	c.GetB().GetCc().Get(1).SetId(7)
	c.GetB().GetCc().Get(1).SetAaa(8)
	c.GetB().GetCc().Get(1).SetBbb(9)
	c.GetB().GetCc().Get(1).SetCcc("10")
	c.GetB().GetCc().Set(1, &nest.TestA{})
	c.GetB().GetCc().Get(0).SetId(11)
	c.GetB().GetCc().Get(0).SetAaa(12)
	c.GetB().GetCc().Get(0).SetBbb(13)
	c.GetB().GetCc().Get(0).SetCcc("14")
	c.GetB().SetDd(gsmodel.NewDMap[string, *nest.TestA]())
	c.GetB().GetDd().Set("c_b_d_1", &nest.TestA{})
	c.GetB().GetDd().Get("c_b_d_1").SetId(15)
	c.GetB().GetDd().Get("c_b_d_1").SetAaa(16)
	c.GetB().GetDd().Get("c_b_d_1").SetBbb(17)
	c.GetB().GetDd().Get("c_b_d_1").SetCcc("18")
	c.SetC(gsmodel.NewDList[*nest.TestB]())
	c.GetC().Append(&nest.TestB{})
	c.GetC().Get(0).SetId(19)
	c.GetC().Get(0).SetAa("20")
	c.GetC().Get(0).SetBb(&nest.TestA{})
	c.GetC().Get(0).GetBb().SetAaa(21)
	c.GetC().Get(0).GetBb().SetBbb(22)
	c.GetC().Get(0).GetBb().SetCcc("23")
	c.GetC().Get(0).SetCc(gsmodel.NewDList[*nest.TestA]())
	c.GetC().Get(0).GetCc().Append(&nest.TestA{})
	c.GetC().Get(0).GetCc().Get(0).SetAaa(24)
	c.GetC().Get(0).GetCc().Get(0).SetBbb(25)
	c.GetC().Get(0).GetCc().Get(0).SetCcc("26")
	c.GetC().Get(0).SetDd(gsmodel.NewDMap[string, *nest.TestA]())
	c.GetC().Get(0).GetDd().Set("c_c_0_d", &nest.TestA{})
	c.GetC().Get(0).GetDd().Get("c_c_0_d").SetId(110)
	c.GetC().Get(0).GetDd().Get("c_c_0_d").SetAaa(111)
	c.GetC().Get(0).GetDd().Get("c_c_0_d").SetBbb(112)
	c.GetC().Get(0).GetDd().Get("c_c_0_d").SetCcc("113")
	c.SetD(gsmodel.NewDMap[string, *nest.TestB]())
	c.GetD().Set("c_d", &nest.TestB{})
	c.GetD().Get("c_d").SetAa("27")
	c.GetD().Get("c_d").SetBb(&nest.TestA{})
	c.GetD().Get("c_d").GetBb().SetAaa(28)
	c.GetD().Get("c_d").GetBb().SetBbb(29)
	c.GetD().Get("c_d").GetBb().SetCcc("30")
	c.GetD().Get("c_d").SetCc(gsmodel.NewDList[*nest.TestA]())
	c.GetD().Get("c_d").GetCc().Append(&nest.TestA{})
	c.GetD().Get("c_d").GetCc().Get(0).SetId(31)
	c.GetD().Get("c_d").GetCc().Get(0).SetAaa(32)
	c.GetD().Get("c_d").GetCc().Get(0).SetBbb(33)
	c.GetD().Get("c_d").GetCc().Get(0).SetCcc("34")
	c.GetD().Get("c_d").SetDd(gsmodel.NewDMap[string, *nest.TestA]())
	c.GetD().Get("c_d").GetDd().Set("c_d_d", &nest.TestA{})
	c.GetD().Get("c_d").GetDd().Get("c_d_d").SetId(35)
	c.GetD().Get("c_d").GetDd().Get("c_d_d").SetAaa(36)
	c.GetD().Get("c_d").GetDd().Get("c_d_d").SetBbb(37)
	c.GetD().Get("c_d").GetDd().Get("c_d_d").SetCcc("38")

	return c
}
func TestMongoLoadSave(t *testing.T) {
	c := getTestC(123)
	fmt.Printf("c:%+v\n", c)
	m1, e := bson.Marshal(c)
	if e != nil {
		panic(e)
	}
	m2 := nest.TestC{}
	e = bson.Unmarshal(m1, &m2)
	if e != nil {
		panic(e)
	}
	fmt.Printf("m2:%+v\n", &m2)

	n1, e := json.Marshal(&c)
	if e != nil {
		panic(e)
	}
	n2 := nest.TestC{}
	e = json.Unmarshal(n1, &n2)
	if e != nil {
		panic(e)
	}
	fmt.Printf("n2:%+v\n", &n2)
	n3, e := n2.Clone()
	if e != nil {
		panic(e)
	}
	n3.SetA("aaa")
	fmt.Printf("n2:%+v\n", &n2)
	fmt.Printf("n3:%+v\n", n3)

	if mongoUrl != "" {
		mongo_helper.Connect(mongoUrl)
		defer mongo_helper.Close()
		col := mongo_helper.GetCol("test", "model")
		filter := bson.M{"_id": c.GetId()}
		_, err := col.ReplaceOne(context.TODO(), filter, c, options.Replace().SetUpsert(true))
		if err != nil {
			t.Error(err)
		}
		zz := nest.TestC{}
		err = col.FindOne(context.TODO(), filter).Decode(&zz)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("n3:%+v\n", &zz)
	}
}
func TestUpdate(t *testing.T) {
	filter := bson.M{"_id": 123}
	update1 := bson.M{"a": "c_a_new", "b.bb.ccc": "1103", "d.c_d.aa": "1027", "d.c_d.bb.bbb": 1029}
	update2 := bson.M{"b.aa": ""}
	if mongoUrl != "" {
		mongo_helper.Connect(mongoUrl)
		defer mongo_helper.Close()
		col := mongo_helper.GetCol("test", "model")
		v, err := col.UpdateOne(context.TODO(), filter, bson.M{"$set": update1, "$unset": update2}, options.Update().SetUpsert(true))
		if err != nil {
			t.Error(err)
		}
		_ = v
		zz := nest.TestC{}
		err = col.FindOne(context.TODO(), filter).Decode(&zz)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("n3:%+v\n", &zz)
	}
}
func TestBuildBson(t *testing.T) {
	filter := bson.M{"_id": 123}
	c := getTestC(123)
	m := bson.M{}
	c.BuildBson(m, "")
	m1 := bson.M{}
	c.BuildBson(m1, "")
	if len(m1) != 0 {
		t.Error("build需要清空dirty")
	}
	c.GetB().GetCc().Get(1).SetId(120)
	c.GetB().GetCc().Get(1).SetAaa(121)
	c.BuildBson(m, "")
	fmt.Printf("n3:%+v\n", m1)
	if mongoUrl != "" {
		mongo_helper.Connect(mongoUrl)
		defer mongo_helper.Close()
		col := mongo_helper.GetCol("test", "model")
		v, err := col.UpdateOne(context.TODO(), filter, m, options.Update().SetUpsert(true))
		if err != nil {
			t.Error(err)
		}
		_ = v
		zz := nest.TestC{}
		err = col.FindOne(context.TODO(), filter).Decode(&zz)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("n3:%+v\n", &zz)
	}
}
