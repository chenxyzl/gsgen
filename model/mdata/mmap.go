package mdata

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

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
}

func NewMMap[K int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string, V any]() *MMap[K, V] {
	ret := &MMap[K, V]{}
	ret.init()
	return ret
}
func (s *MMap[K, V]) init() {
	s.data = make(map[K]V)
	s.dirty = make(map[K]bool)
}

// Len 长度
func (s *MMap[K, V]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

// Reset 重置清空list
func (s *MMap[K, V]) Reset() {
	if s == nil {
		return
	}
	s.data = make(map[K]V)
	s.updateDirtyAll()
}

// Get 设置值
func (s *MMap[K, V]) Get(k K) V {
	if s == nil {
		panic("map is nil")
	}
	return s.data[k]
}

// Set 设置新值
func (s *MMap[K, V]) Set(k K, v V) {
	if s == nil {
		panic("map is nil")
	}
	//
	checkSetParent(v, k, s.updateDirty)
	s.data[k] = v
	s.updateDirty(k)
}

// Remove 删除 注:因为删除不太好处理list对应的bson的更新,所以这里用了DirtyAll
func (s *MMap[K, V]) Remove(k K) {
	if s == nil {
		panic("map is nil")
	}
	if _, ok := s.data[k]; !ok {
		return
	}
	delete(s.data, k)
	s.updateDirty(k)
}

// Range 遍历
func (s *MMap[K, V]) Range(f func(K, V) bool) {
	if s == nil {
		panic("map is nil")
	}
	if f == nil {
		return
	}
	for k, v := range s.data {
		if _continue := f(k, v); !_continue {
			break
		}
	}
}

// SetParent 设置父节点
func (s *MMap[K, V]) SetParent(idx any, dirtyParentFunc DirtyParentFunc) {
	if s == nil {
		return
	}
	if s.dirtyParent != nil {
		panic("model被重复设置了父节点,请先从老节点移除")
	}
	s.inParentDirtyIdx = idx
	s.dirtyParent = dirtyParentFunc
}

// IsDirty 是否为脏
func (s *MMap[K, V]) IsDirty() bool {
	if s.dirtyAll {
		return true
	}
	return len(s.dirty) > 0
}

// CleanDirty 清楚脏标记
func (s *MMap[K, V]) CleanDirty(withChildren bool) {
	if s == nil {
		return
	}
	if withChildren {
		if s.dirtyAll {
			var v V
			if _, ok := (any(v)).(IDirtyModel); ok {
				s.Range(func(k K, v V) bool {
					(any(v)).(IDirtyModel).CleanDirty(withChildren)
					return true
				})
			}
		} else {
			var v V
			if _, ok := (any(v)).(IDirtyModel); ok {
				for nk := range s.dirty {
					(any(s.Get(nk))).(IDirtyModel).CleanDirty(withChildren)
				}
			}

		}
	}
	s.dirtyAll = false
	clear(s.dirty)
}

// updateDirty 更新藏标记
func (s *MMap[K, V]) updateDirty(tk any) {
	k := tk.(K)
	//如果已经allDirty了就不用管了
	if s.dirtyAll || s.dirty[k] {
		return
	}
	s.dirty[k] = true
	if s.dirtyParent != nil {
		s.dirtyParent.Invoke(s.inParentDirtyIdx)
	}
}

// updateDirtyAll 更新整个对象为脏
func (s *MMap[K, V]) updateDirtyAll() {
	if s.dirtyAll {
		return
	}
	s.dirtyAll = true
	if s.dirtyParent != nil {
		s.dirtyParent.Invoke(s.inParentDirtyIdx)
	}
}

// String toString
func (s *MMap[K, V]) String() string {
	return fmt.Sprintf("%v", s.data)
}

// MarshalJSON json序列化
func (s *MMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

// UnmarshalJSON json反序列化
func (s *MMap[K, V]) UnmarshalJSON(data []byte) error {
	var m map[K]V
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	s.init()
	for k, v := range m {
		s.Set(k, v)
	}
	return nil
}

// MarshalBSON bson序列化
func (s *MMap[K, V]) MarshalBSON() ([]byte, error) {
	r, r1, r2 := bson.MarshalValue(s.data)
	_ = r
	return r1, r2
}

// UnmarshalBSON bson反序列化
func (s *MMap[K, V]) UnmarshalBSON(data []byte) error {
	var m map[K]V
	if err := bson.UnmarshalValue(bson.TypeEmbeddedDocument, data, &m); err != nil {
		return err
	}
	s.init()
	for k, v := range m {
		s.Set(k, v)
	}
	return nil
}

// BuildDirty bson的增量更新
func (s *MMap[K, V]) BuildDirty(m bson.M, preKey string) {
	if len(s.dirty) == 0 && !s.dirtyAll {
		return
	}
	if s.dirtyAll {
		AddSetDirtyM(m, preKey, s)
	} else {
		for k := range s.dirty {
			AddSetDirtyM(m, MakeBsonKey(fmt.Sprintf("%v", k), preKey), s.data[k])
		}
	}
	s.CleanDirty(false)
	return
}
