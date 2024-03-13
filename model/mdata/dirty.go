package mdata

import (
	"math"
)

const DirtyAll = math.MaxUint64

type DirtyParentFunc func(dirtyIdx uint64)

func (f DirtyParentFunc) Invoke(dirtyIdx uint64) {
	if f != nil {
		f(dirtyIdx)
	}
}

// DirtyModel ------------------dirtyModel脏标记
type DirtyModel struct {
	dirty        uint64
	selfDirtyIdx uint64
	dirtyParent  DirtyParentFunc
}

func (this *DirtyModel) SetSelfDirtyIdx(idx uint64, dirtyParentFunc DirtyParentFunc) {
	this.selfDirtyIdx = idx
	this.dirtyParent = dirtyParentFunc
}
func (this *DirtyModel) IsDirty() bool {
	return this.dirty > 0
}
func (this *DirtyModel) IsDirtyAll() bool {
	return this.dirty&DirtyAll == DirtyAll
}
func (this *DirtyModel) DirtyAll() {
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
		this.dirtyParent.Invoke(this.selfDirtyIdx)
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
