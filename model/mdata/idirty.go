package mdata

// test
type testDirtyModel struct{ DirtyModel }

// check
var _ IDirtyModel = (*testDirtyModel)(nil)
var _ IDirtyModel = (*MList[int])(nil)
var _ IDirtyModel = (*MMap[int, int])(nil)

type IDirtyModel interface {
	SetParent(idx any, dirtyParentFunc DirtyParentFunc)
	IsDirty() bool
	CleanDirty()
}

func CheckCallDirty(v any, idx any, dirtyParentFunc DirtyParentFunc) {
	if dirty, ok := v.(IDirtyModel); ok {
		dirty.SetParent(idx, dirtyParentFunc)
	}
}
