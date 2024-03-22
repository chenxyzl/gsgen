package mdata

import (
	"math"
)

const DirtyAll = math.MaxUint64

type DirtyParentFunc func(dirtyIdx any)

func (f DirtyParentFunc) Invoke(dirtyIdx any) {
	if f != nil {
		f(dirtyIdx)
	}
}

// DirtyModel ------------------dirtyModel脏标记
type DirtyModel struct {
	dirty            uint64
	inParentDirtyIdx any
	dirtyParent      DirtyParentFunc
}

func (this *DirtyModel) SetParent(idx any, dirtyParentFunc DirtyParentFunc) {
	if this == nil {
		return
	}
	if this.dirtyParent != nil {
		panic("model被重复设置了父节点,请先从老节点移除")
	}
	this.inParentDirtyIdx = idx
	this.dirtyParent = dirtyParentFunc
}
func (this *DirtyModel) IsDirty() bool {
	return this.dirty > 0
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
func (this *DirtyModel) UpdateDirty(tn any) {
	n := uint64(tn.(int))
	if this.dirty&n == n {
		return
	}
	this.dirty |= n
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.inParentDirtyIdx)
	}
}
