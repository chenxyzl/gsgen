package bson

import "github.com/chenxyzl/gsgen/gsmodel"

type TestA struct {
	gsmodel.DirtyModel `bson:"-"`
	cc                 *gsmodel.DList[int]           `bson:"cc"`
	dd                 *gsmodel.DMap[string, *TestA] `bson:"dd"`
}
