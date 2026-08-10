package bson

import "github.com/chenxyzl/gsgen/gsmodel"

type TestA struct {
	gsmodel.DirtyModel `bson:"-"`
	cc                 *gsmodel.DList[int]           `bson:"cc"`
	dd                 *gsmodel.DMap[string, *TestA] `bson:"dd"`
	ignoreMe           int64                         `bson:"-"` //测试:命名字段的bson:"-"忽略,仍生成getter/setter/dirty,但不参与bson
}
