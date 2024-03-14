package mdata

// test
type testDirtyModel struct{ DirtyModel }

// check
var _ IDirtyModel[uint64] = &testDirtyModel{}
var _ IDirtyModel[uint64] = &MList[int]{}
var _ IDirtyModel[int] = &MMap[int, int]{}

type IDirtyModel[T comparable] interface {
	SetSelfDirtyIdx(idx T, dirtyParentFunc DirtyParentFunc[T])
	IsDirty() bool
	IsDirtyAll() bool
	UpdateDirty(n T)
	UpdateDirtyAll()
	CleanDirty()
}

func CheckCallDirty[T comparable](v any, idx T, dirtyParentFunc DirtyParentFunc[T]) {
	if dirty, ok := v.(IDirtyModel[T]); ok {
		dirty.SetSelfDirtyIdx(idx, dirtyParentFunc)
	}
}
