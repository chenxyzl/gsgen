package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"gotest/model/mdata"
)

func (c *TestA) MarshalBSON() ([]byte, error) {
	doc := bson.M{
		"_id": c.id,
		"a":   c.a,
		"b":   c.b,
	}
	return bson.Marshal(doc)
}

func (c *TestA) UnmarshalBSON(data []byte) error {
	doc := struct {
		Id uint64 `bson:"_id"`
		A  int64  `bson:"a"`
		B  int32  `bson:"b"` //
	}{}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	c.SetId(doc.Id)
	c.SetA(doc.A)
	c.SetB(doc.B)
	return nil
}
func (m *TestB) MarshalBSON() ([]byte, error) {
	a := struct {
		Id uint64                      `bson:"_id"`
		M  string                      `bson:"m"`
		N  *TestA                      `bson:"n"` // 内嵌结构体类型 A
		C  *mdata.MList[*TestA]        `bson:"c"`
		D  *mdata.MMap[string, *TestA] `bson:"d"`
	}{Id: m.id, M: m.m, N: m.n, C: m.c, D: m.d}
	return bson.Marshal(&a)
}
func (m *TestB) UnmarshalBSON(data []byte) error {
	doc := struct {
		Id uint64                      `bson:"_id"`
		M  string                      `bson:"m"`
		N  *TestA                      `bson:"n"` //
		C  *mdata.MList[*TestA]        `bson:"c"`
		D  *mdata.MMap[string, *TestA] `bson:"d"`
	}{}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	m.SetId(doc.Id)
	m.SetM(doc.M)
	m.SetN(doc.N)
	m.SetC(doc.C)
	m.SetD(doc.D)
	return nil
}
func (m *TestC) MarshalBSON() ([]byte, error) {
	doc := bson.M{
		"_id": m.id,
		"x":   m.x,
		"y":   m.y,
	}
	return bson.Marshal(doc)
}
func (m *TestC) UnmarshalBSON(data []byte) error {
	doc := struct {
		Id uint64 `bson:"_id"`
		X  string `bson:"x"`
		Y  *TestB `bson:"y"` //
	}{}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	m.SetId(doc.Id)
	m.SetX(doc.X)
	m.SetY(doc.Y)
	return nil
}
