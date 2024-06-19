package with_ignore

import (
	"github.com/chenxyzl/gsgen/example/common"
	"github.com/chenxyzl/gsgen/gsmodel"
)

type TestA struct {
	gsmodel.DirtyModel `bson:"-"`
	ig                 *common.Common `bson:"ig"`
}
