package getter

import (
	"github.com/chenxyzl/gsgen/gsmodel"
)

type TestA struct {
	cc *gsmodel.AList[int]           `bson:"cc"`
	dd *gsmodel.AMap[string, *TestA] `bson:"dd"`
}
