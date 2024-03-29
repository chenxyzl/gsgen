package mdata

// test
type testDirtyModel struct{ DirtyModel }

// check
var _ IDirtyModel = (*testDirtyModel)(nil)
var _ IDirtyModel = (*MList[int])(nil)
var _ IDirtyModel = (*MMap[int, int])(nil)

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
