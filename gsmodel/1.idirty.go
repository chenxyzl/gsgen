package gsmodel

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

// test
type testDirtyModel struct{ DirtyModel }

// check
var _ IDirtyModel = (*testDirtyModel)(nil)
var _ IDirtyModel = (*DList[int])(nil)
var _ IDirtyModel = (*DMap[int, int])(nil)

// IDirtyModel model接口
type IDirtyModel interface {
	SetParent(idx any, dirtyParentFunc DirtyParentFunc)
	IsDirty() bool
	CleanDirty()
}

// checkSetParent 设置对象的父节点
func checkSetParent(v any, idx any, dirtyParentFunc DirtyParentFunc) {
	if dirty, ok := v.(IDirtyModel); ok {
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
