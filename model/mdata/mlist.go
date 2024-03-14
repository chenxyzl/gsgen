package mdata

import (
	"fmt"
)

// MList ----------------------MList-----------------------
type MList[T any] struct {
	data []T `bson:"data"`
	//
	dirty        map[uint64]bool
	dirtyAll     bool
	selfDirtyIdx uint64
	dirtyParent  DirtyParentFunc[uint64]
}

func NewList[T any]() *MList[T] {
	return &MList[T]{data: make([]T, 0), dirty: make(map[uint64]bool)}
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
	this.UpdateDirtyAll()
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
	if idx >= l {
		panic(fmt.Sprintf("MList set idx out of range, len:%d|idx:%d", l, idx))
	}
	//
	CheckCallDirty(v, idx, this.UpdateDirty)
	this.data[idx] = v
	this.UpdateDirty(idx)
}

// Append 追加
func (this *MList[T]) Append(vs ...T) {
	if this == nil {
		panic("data is nil")
	}
	for _, v := range vs {
		idx := uint64(this.Len())
		//
		CheckCallDirty(v, idx, this.UpdateDirty)
		this.data = append(this.data, v)
		this.UpdateDirty(idx)
	}
}

// Remove 删除 注:因为删除不太好处理list对应的mongo的更新,所以这里用了DirtyAll
func (this *MList[T]) Remove(idx int) {
	if this == nil {
		panic("data is nil")
	}
	l := this.Len()
	if idx >= l {
		panic(fmt.Sprintf("MList remove idx out of range, len:%d|idx:%d", l, idx))
	}
	this.data = append(this.data[0:idx], this.data[idx+1:]...)
	this.UpdateDirtyAll()
}

// Range 遍历
func (this *MList[T]) Range(f func(idx int, v T)) {
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

func (this *MList[T]) SetSelfDirtyIdx(idx uint64, dirtyParentFunc DirtyParentFunc[uint64]) {
	if this.dirtyParent != nil {
		panic("model被重复设置了父节点,请先从老节点移除")
	}
	this.selfDirtyIdx = idx
	this.dirtyParent = dirtyParentFunc
}

func (this *MList[T]) IsDirty() bool {
	return len(this.dirty) > 0
}

func (this *MList[T]) IsDirtyAll() bool {
	return this.dirtyAll
}

func (this *MList[T]) UpdateDirty(n uint64) {
	//如果已经allDirty了就不用管了
	if this.dirtyAll || this.dirty[n] {
		return
	}
	this.dirty[n] = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.selfDirtyIdx)
	}
}
func (this *MList[T]) UpdateDirtyAll() {
	if this.dirtyAll {
		return
	}
	this.dirtyAll = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.selfDirtyIdx)
	}
}
func (this *MList[T]) CleanDirty() {
	if this == nil {
		return
	}
	var v T //todo 类型不一定是uint64
	if _, ok := (any(v)).(IDirtyModel[uint64]); ok {
		this.Range(func(idx int, v T) {
			(any(v)).(IDirtyModel[uint64]).CleanDirty()
		})
	}
	clear(this.dirty)
}
