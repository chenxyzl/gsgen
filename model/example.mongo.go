package model

import (
	"go.mongodb.org/mongo-driver/bson"
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
	var doc bson.M
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	c.SetId(uint64(doc["_id"].(int64)))
	c.SetA(doc["a"].(int64))
	c.SetB(doc["b"].(int32))
	return nil
}
func (m *TestB) MarshalBSON() ([]byte, error) {
	doc := bson.M{
		"_id": m.id,
		"m":   m.m,
		"n":   m.n,
	}
	return bson.Marshal(doc)
}
func (m *TestB) UnmarshalBSON(data []byte) error {
	var doc bson.M
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	m.SetId(uint64(doc["_id"].(int64)))
	m.SetM(doc["m"].(string))
	//The problem is here
	//repeated Marshal and Unmarshal
	//convert primitive.M to struct TestA
	//todo 需要更优雅的方法
	if dat, err := bson.Marshal(doc["n"]); err != nil {
		return err
	} else {
		var n TestA
		if err := bson.Unmarshal(dat, &n); err != nil {
			return err
		}
		m.SetN(&n)
	}
	return nil
}
