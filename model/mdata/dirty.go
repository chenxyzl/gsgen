package mdata

import (
	"math"
)

const DirtyAll = math.MaxUint64

type DirtyParentFunc[T comparable] func(dirtyIdx T)

func (f DirtyParentFunc[T]) Invoke(dirtyIdx T) {
	if f != nil {
		f(dirtyIdx)
	}
}

// DirtyModel ------------------dirtyModel脏标记
type DirtyModel struct {
	dirty            uint64
	inParentDirtyIdx uint64
	dirtyParent      DirtyParentFunc[uint64]
}

func (this *DirtyModel) SetParent(idx uint64, dirtyParentFunc DirtyParentFunc[uint64]) {
	if this.dirtyParent != nil {
		panic("model被重复设置了父节点,请先从老节点移除")
	}
	this.inParentDirtyIdx = idx
	this.dirtyParent = dirtyParentFunc
}
func (this *DirtyModel) IsDirty() bool {
	return this.dirty > 0
}
func (this *DirtyModel) IsDirtyAll() bool {
	return this.dirty&DirtyAll == DirtyAll
}
func (this *DirtyModel) UpdateDirtyAll() {
	if this == nil {
		return
	}
	this.UpdateDirty(DirtyAll)
}
func (this *DirtyModel) UpdateDirty(n uint64) {
	if this.dirty&n == n {
		return
	}
	this.dirty |= n
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.inParentDirtyIdx)
	}
}
func (this *DirtyModel) CleanDirty() {
	if this == nil {
		return
	}
	d := this.dirty
	if d == 0 {
		return
	} else {
		this.dirty = 0
	}
}
