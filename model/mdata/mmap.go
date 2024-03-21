package mdata

import "go.mongodb.org/mongo-driver/bson"

// MMap ----------------------------------MMap-------------------------------------
// MMap map的包装
// @K key的类型
// @V value的类型
type MMap[K int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string, V any] struct {
	data map[K]V `bson:"map"`
	//
	dirty            map[K]bool
	dirtyAll         bool
	inParentDirtyIdx any
	dirtyParent      DirtyParentFunc
	isNumKey         bool
}

func NewMMap[K int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string, V any]() *MMap[K, V] {
	ret := &MMap[K, V]{}
	ret.init()
	return ret
}
func (this *MMap[K, V]) init() {
	this.data = make(map[K]V)
	this.dirty = make(map[K]bool)
	var k K
	this.isNumKey = isNum(k)
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
	this.updateDirtyAll()
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
	//
	CheckCallDirty(v, k, this.updateDirty)
	this.data[k] = v
	this.updateDirty(k)
}

// Remove 删除 注:因为删除不太好处理list对应的mongo的更新,所以这里用了DirtyAll
func (this *MMap[K, V]) Remove(k K) {
	if this == nil {
		panic("map is nil")
	}
	if _, ok := this.data[k]; !ok {
		return
	}
	delete(this.data, k)
	this.updateDirty(k)
}

// Range 遍历
func (this *MMap[K, V]) Range(f func(K, V) bool) {
	if this == nil {
		panic("map is nil")
	}
	if f == nil {
		return
	}
	for k, v := range this.data {
		if _continue := f(k, v); !_continue {
			break
		}
	}
}

func (this *MMap[K, V]) SetParent(idx any, dirtyParentFunc DirtyParentFunc) {
	if this.dirtyParent != nil {
		panic("model被重复设置了父节点,请先从老节点移除")
	}
	this.inParentDirtyIdx = idx
	this.dirtyParent = dirtyParentFunc
}

func (this *MMap[K, V]) IsDirty() bool {
	if this.dirtyAll {
		return true
	}
	if this.isNumKey {
		return len(this.dirty) > 0
	} else {
		isDirty := false
		this.Range(func(k K, v V) bool {
			if mod, ok := any(v).(IDirtyModel); ok {
				if mod.IsDirty() {
					isDirty = true
					return false
				} else {
					return true
				}
			}
			return true
		})
		return isDirty
	}
}

func (this *MMap[K, V]) CleanDirty() {
	if this == nil {
		return
	}
	if this.dirtyAll {
		var v V
		if _, ok := (any(v)).(IDirtyModel); ok {
			this.Range(func(k K, v V) bool {
				(any(v)).(IDirtyModel).CleanDirty()
				return true
			})
		}
	} else {
		var v V
		if _, ok := (any(v)).(IDirtyModel); ok {
			for nk := range this.dirty {
				(any(this.Get(nk))).(IDirtyModel).CleanDirty()
			}
		}

	}
	this.dirtyAll = false
	clear(this.dirty)
}

// updateDirty 更新藏标记
func (this *MMap[K, V]) updateDirty(tk any) {
	k := tk.(K)
	//如果已经allDirty了就不用管了
	if this.dirtyAll || this.dirty[k] {
		return
	}
	this.dirty[k] = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.inParentDirtyIdx)
	}
}
func (this *MMap[K, V]) updateDirtyAll() {
	if this.dirtyAll {
		return
	}
	this.dirtyAll = true
	if this.dirtyParent != nil {
		this.dirtyParent.Invoke(this.inParentDirtyIdx)
	}
}

func (this *MMap[K, V]) MarshalBSON() ([]byte, error) {
	r, r1, r2 := bson.MarshalValue(this.data)
	_ = r
	return r1, r2
}
func (this *MMap[K, V]) UnmarshalBSON(data []byte) error {
	var m map[K]V
	if err := bson.UnmarshalValue(bson.TypeEmbeddedDocument, data, &m); err != nil {
		return err
	}
	this.init()
	for k, v := range m {
		this.Set(k, v)
	}
	return nil
}
