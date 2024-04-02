package model

import (
	"gen_tools/model/mdata"
)

type TestA struct {
	id               uint64 `bson:"_id"`
	aaa              int64  `bson:"aaa"`
	bbb              int32  `bson:"bbb"`
	ccc              string `bson:"ccc"`
	mdata.DirtyModel `bson:"-"`
}
type TestB struct {
	id               uint64                      `bson:"_id"`
	aa               string                      `bson:"aa"`
	bb               *TestA                      `bson:"bb"` // 内嵌结构体类型 A
	cc               *mdata.MList[*TestA]        `bson:"cc"`
	dd               *mdata.MMap[string, *TestA] `bson:"dd"`
	mdata.DirtyModel `bson:"-"`
}
type TestC struct {
	id               uint64                      `bson:"_id"`
	a                string                      `bson:"a"`
	b                *TestB                      `bson:"b"` // 内嵌结构体类型 B
	c                *mdata.MList[*TestB]        `bson:"c"`
	d                *mdata.MMap[string, *TestB] `bson:"d"`
	mdata.DirtyModel `bson:"-"`
}
