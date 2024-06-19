package gsmodel

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

// test
type testDirtyModel struct{ DirtyModel }

// check
var _ iDirtyModel = (*testDirtyModel)(nil)
var _ iDirtyModel = (*DList[int])(nil)
var _ iDirtyModel = (*DMap[int, int])(nil)

// iDirtyModel model接口
type iDirtyModel interface {
	SetParent(idx any, dirtyParentFunc dirtyParentFunc)
	IsDirty() bool
	CleanDirty()
}

// checkSetParent 设置对象的父节点
func checkSetParent(v any, idx any, dirtyParentFunc dirtyParentFunc) {
	if dirty, ok := v.(iDirtyModel); ok {
		dirty.SetParent(idx, dirtyParentFunc)
	}
}

// MakeBsonKey 构造bson.key
func MakeBsonKey(key any, preKey string) string {
	if preKey == "" {
		return fmt.Sprintf("%v", key)
	}
	return fmt.Sprintf("%v.%v", preKey, key)
}

func AddSetDirtyM(m bson.M, k string, v any) {
	if _, ok := m["$set"]; !ok {
		m["$set"] = bson.M{}
	}
	m["$set"].(bson.M)[k] = v
}

func AddUnsetDirtyM(m bson.M, k string) {
	if _, ok := m["$unset"]; !ok {
		m["$unset"] = bson.M{}
	}
	m["$unset"].(bson.M)[k] = ""
}
