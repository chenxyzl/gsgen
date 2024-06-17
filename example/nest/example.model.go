package nest

import (
	"github.com/chenxyzl/gsgen/gsmodel"
)

type TestA struct {
	id                 uint64 `bson:"_id"`
	aaa                int64  `bson:"aaa"`
	bbb                int32  `bson:"bbb"`
	ccc                string `bson:"ccc"`
	gsmodel.DirtyModel `bson:"-"`
}
type TestB struct {
	id                 uint64                        `bson:"_id"`
	aa                 string                        `bson:"aa"`
	bb                 *TestA                        `bson:"bb"` // 内嵌结构体类型 A
	cc                 *gsmodel.DList[*TestA]        `bson:"cc"`
	dd                 *gsmodel.DMap[string, *TestA] `bson:"dd"`
	gsmodel.DirtyModel `bson:"-"`
}
type TestC struct {
	id                 uint64                        `bson:"_id"`
	a                  string                        `bson:"a"`
	b                  *TestB                        `bson:"b"` // 内嵌结构体类型 B
	c                  *gsmodel.DList[*TestB]        `bson:"c"`
	d                  *gsmodel.DMap[string, *TestB] `bson:"d"`
	gsmodel.DirtyModel `bson:"-"`
}
