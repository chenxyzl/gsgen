package model

import "gotest/model/mdata"

type TestA struct {
	id               uint64 `bson:"_id"`
	a                int64  `bson:"a"`
	b                int32  `bson:"b"`
	mdata.DirtyModel `bson:"-"`
}

type TestB struct {
	id               uint64 `bson:"_id"`
	m                string `bson:"n"`
	n                *TestA `bson:"n"` // 内嵌结构体类型 A
	mdata.DirtyModel `bson:"-"`
}
