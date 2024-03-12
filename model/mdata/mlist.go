package mdata

import (
	"fmt"
)

// MList ----------------------MList-----------------------
type MList[T any] struct {
	data []T `bson:"data"`
	//
	dirty        map[uint64]bool
	selfDirtyIdx uint64
	dirtyParent  DirtyParentFunc
}

func NewList[T any]() *MList[T] {
	return &MList[T]{data: make([]T, 0), dirty: make(map[uint64]bool)}
}

// SetSelfDirtyIdx 设置父节点
func (this *MList[T]) SetSelfDirtyIdx(idx uint64, dirtyParentFunc DirtyParentFunc) {
	this.selfDirtyIdx = idx
	this.dirtyParent = dirtyParentFunc
}

// Len 长度
func (this *MList[T]) Len() int {
	if this == nil {
		return 0
	}
	return len(this.data)
}

// Reset 重置清空list
func (this *MList[T]) Reset() {
	if this == nil {
		return
	}
	this.data = make([]T, 0)
	this.updateDirty(DirtyAll)
}

// Get 设置值
func (this *MList[T]) Get(idx int) T {
	if this == nil {
		panic("MList is nil")
	}
	l := this.Len()
	if idx >= l {
		panic(fmt.Sprintf("MList get idx out of range, len:%d|idx:%d", l, idx))
	}
	return this.data[idx]
}

// Set 设置新值
func (this *MList[T]) Set(idx uint64, v T) {
	if this == nil {
		panic("data is nil")
	}
	l := uint64(this.Len())
	if l >= idx {
		panic(fmt.Sprintf("MList set idx out of range, len:%d|idx:%d", l, idx))
	}
	//todo v的类型如果是非基本类型则需要设置dirtyParent
	this.data[idx] = v
	//
	this.updateDirty(idx)
}

// Append 追加
func (this *MList[T]) Append(vs ...T) {
	if this == nil {
		panic("data is nil")
	}
	for _, v := range vs {
		//todo v的类型如果是非基本类型则需要设置dirtyParent
		//
		this.data = append(this.data, v)
		//
		this.updateDirty(uint64(this.Len()))
	}
}

// Remove 删除 注:因为删除不太好处理list对应的mongo的更新,所以这里用了DirtyAll
func (this *MList[T]) Remove(idx int) {
	if this == nil {
		panic("data is nil")
	}
	l := this.Len()
	if l >= idx {
		panic(fmt.Sprintf("MList remove idx out of range, len:%d|idx:%d", l, idx))
	}
	this.data = append(this.data[0:idx], this.data[idx+1:]...)
	this.updateDirty(DirtyAll)
}

// Range 遍历
func (this *MList[T]) Range(f func(idx int, value T)) {
	if this == nil {
		panic("MList is nil")
	}
	if f == nil {
		return
	}
	for idx, v := range this.data {
		f(idx, v)
	}
}

// updateDirty 更新藏标记
func (this *MList[T]) updateDirty(n uint64) {
	if this.dirty[n] {
		return
	}
	this.dirty[n] = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.selfDirtyIdx)
	}
}
