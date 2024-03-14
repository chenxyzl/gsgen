package mdata

// MMap ----------------------------------MMap-------------------------------------
type MMap[K comparable, v any] struct {
	data map[K]v `bson:"map"`
	//
	dirty        map[K]bool
	dirtyAll     bool
	selfDirtyIdx K
	dirtyParent  DirtyParentFunc[K]
}

func NewMMap[K comparable, V any]() *MMap[K, V] {
	return &MMap[K, V]{data: make(map[K]V), dirty: make(map[K]bool)}
}

// Len 长度
func (this *MMap[K, V]) Len() int {
	if this == nil {
		return 0
	}
	return len(this.data)
}

// Reset 重置清空list
func (this *MMap[K, V]) Reset() {
	if this == nil {
		return
	}
	this.data = make(map[K]V)
	this.UpdateDirtyAll()
}

// Get 设置值
func (this *MMap[K, V]) Get(k K) V {
	if this == nil {
		panic("map is nil")
	}
	return this.data[k]
}

// Set 设置新值
func (this *MMap[K, V]) Set(k K, v V) {
	if this == nil {
		panic("map is nil")
	}
	//todo v的类型如果是非基本类型则需要设置dirtyParent
	//CheckCallDirty(v, 0, this.UpdateDirty)
	this.data[k] = v
	//
	this.UpdateDirty(k)
}

// Remove 删除 注:因为删除不太好处理list对应的mongo的更新,所以这里用了DirtyAll
func (this *MMap[K, V]) Delete(k K) {
	if this == nil {
		panic("map is nil")
	}
	if _, ok := this.data[k]; !ok {
		return
	}
	delete(this.data, k)
	this.UpdateDirty(k)
}

// Range 遍历
func (this *MMap[K, V]) Range(f func(K, V)) {
	if this == nil {
		panic("map is nil")
	}
	if f == nil {
		return
	}
	for k, v := range this.data {
		f(k, v)
	}
}

func (this MMap[K, v]) SetSelfDirtyIdx(idx K, dirtyParentFunc DirtyParentFunc[K]) {
	if this.dirtyParent != nil {
		panic("model被重复设置了父节点,请先从老节点移除")
	}
	this.selfDirtyIdx = idx
	this.dirtyParent = dirtyParentFunc
}

func (this MMap[K, v]) IsDirty() bool {
	return len(this.dirty) > 0 || this.dirtyAll
}

func (this MMap[K, v]) IsDirtyAll() bool {
	return this.dirtyAll
}

func (this MMap[K, v]) DirtyAll() {
	this.dirtyAll = true
}

// updateDirty 更新藏标记
func (this *MMap[K, V]) UpdateDirty(k K) {
	//如果已经allDirty了就不用管了
	if this.dirtyAll || this.dirty[k] {
		return
	}
	this.dirty[k] = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.selfDirtyIdx)
	}
}
func (this *MMap[K, V]) UpdateDirtyAll() {
	if this.dirtyAll {
		return
	}
	this.dirtyAll = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.selfDirtyIdx)
	}
}
func (this *MMap[K, V]) CleanDirty() {
	if this == nil {
		return
	}
	var v V //类型不一定是uint64
	if _, ok := (any(v)).(IDirtyModel[uint64]); ok {
		this.Range(func(k K, v V) {
			(any(v)).(IDirtyModel[uint64]).CleanDirty()
		})
	}
	clear(this.dirty)
}
