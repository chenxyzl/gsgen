package mdata

// MMap ----------------------------------MMap-------------------------------------
// MMap map的包装
// @K key的类型
// @V value的类型
// @D 在父节点的脏标记类型
type MMap[K comparable, V any, D comparable] struct {
	data map[K]V `bson:"map"`
	//
	dirty            map[K]bool
	dirtyAll         bool
	inParentDirtyIdx D
	dirtyParent      DirtyParentFunc[D]
}

func NewMMap[K comparable, V any, D comparable]() *MMap[K, V, D] {
	return &MMap[K, V, D]{data: make(map[K]V), dirty: make(map[K]bool)}
}

// Len 长度
func (this *MMap[K, V, D]) Len() int {
	if this == nil {
		return 0
	}
	return len(this.data)
}

// Reset 重置清空list
func (this *MMap[K, V, D]) Reset() {
	if this == nil {
		return
	}
	this.data = make(map[K]V)
	this.UpdateDirtyAll()
}

// Get 设置值
func (this *MMap[K, V, D]) Get(k K) V {
	if this == nil {
		panic("map is nil")
	}
	return this.data[k]
}

// Set 设置新值
func (this *MMap[K, V, D]) Set(k K, v V) {
	if this == nil {
		panic("map is nil")
	}
	//todo v的类型如果是非基本类型则需要设置dirtyParent
	CheckCallDirty[K](v, k, this.UpdateDirty)
	this.data[k] = v
	//
	this.UpdateDirty(k)
}

// Remove 删除 注:因为删除不太好处理list对应的mongo的更新,所以这里用了DirtyAll
func (this *MMap[K, V, D]) Delete(k K) {
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
func (this *MMap[K, V, D]) Range(f func(K, V)) {
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

func (this MMap[K, V, D]) SetParent(idx D, dirtyParentFunc DirtyParentFunc[D]) {
	if this.dirtyParent != nil {
		panic("model被重复设置了父节点,请先从老节点移除")
	}
	this.inParentDirtyIdx = idx
	this.dirtyParent = dirtyParentFunc
}

func (this MMap[K, V, D]) IsDirty() bool {
	return len(this.dirty) > 0 || this.dirtyAll
}

func (this MMap[K, V, D]) IsDirtyAll() bool {
	return this.dirtyAll
}

func (this MMap[K, V, D]) DirtyAll() {
	this.dirtyAll = true
}

// updateDirty 更新藏标记
func (this *MMap[K, V, D]) UpdateDirty(k K) {
	//如果已经allDirty了就不用管了
	if this.dirtyAll || this.dirty[k] {
		return
	}
	this.dirty[k] = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.inParentDirtyIdx)
	}
}
func (this *MMap[K, V, D]) UpdateDirtyAll() {
	if this.dirtyAll {
		return
	}
	this.dirtyAll = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.inParentDirtyIdx)
	}
}
func (this *MMap[K, V, D]) CleanDirty() {
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
