package gsmodel

import (
	"encoding/json"
	"fmt"
	"maps"

	"go.mongodb.org/mongo-driver/bson"
)

// DMap ----------------------------------DMap-------------------------------------
// DMap map的包装
// @K key的类型
// @V value的类型
type DMap[K int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string, V any] struct {
	data map[K]V `bson:"map"`
	//
	dirty            map[K]bool //被新增/修改的key,生成$set
	removed          map[K]bool //被删除的key,生成$unset
	dirtyAll         bool
	inParentDirtyIdx any
	dirtyParent      dirtyParentFunc
}

func NewDMap[K int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string, V any]() *DMap[K, V] {
	ret := &DMap[K, V]{}
	ret.init()
	return ret
}
func (s *DMap[K, V]) init() {
	s.data = make(map[K]V)
	s.dirty = make(map[K]bool)
	s.removed = make(map[K]bool)
}

// Len 长度
func (s *DMap[K, V]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

// Clean 重置清空list
func (s *DMap[K, V]) Clean() {
	if s == nil || len(s.data) == 0 {
		return
	}
	s.data = make(map[K]V)
	s.updateDirtyAll()
}

// Get 设置值
func (s *DMap[K, V]) Get(k K) V {
	if s == nil {
		panic("map is nil")
	}
	return s.data[k]
}

// Set 设置新值
func (s *DMap[K, V]) Set(k K, v V) {
	if s == nil {
		panic("map is nil")
	}
	//
	checkSetParent(v, k, s.updateDirty)
	s.data[k] = v
	delete(s.removed, k) //重新赋值,取消删除标记
	s.updateDirty(k)
}

// Remove 删除
func (s *DMap[K, V]) Remove(k K) {
	if s == nil {
		panic("map is nil")
	}
	if _, ok := s.data[k]; !ok {
		return
	}
	delete(s.data, k)
	delete(s.dirty, k) //取消$set标记
	s.updateRemoved(k) //标记为$unset
}

// Range 遍历
func (s *DMap[K, V]) Range(f func(K, V) bool) {
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
func (s *DMap[K, V]) SetParent(idx any, dirtyParentFunc dirtyParentFunc) {
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
func (s *DMap[K, V]) IsDirty() bool {
	if s == nil {
		return false
	}
	return s.dirtyAll || len(s.dirty) > 0 || len(s.removed) > 0
}

// CleanDirty 清楚脏标记
func (s *DMap[K, V]) CleanDirty() {
	if s == nil {
		return
	}
	//仅当value本身是脏模型时才需要下钻清理(基本类型的V不实现iDirtyModel)
	var zero V
	if _, ok := any(zero).(iDirtyModel); ok {
		if s.dirtyAll {
			s.Range(func(k K, v V) bool {
				if dm, ok := any(v).(iDirtyModel); ok {
					dm.CleanDirty()
				}
				return true
			})
		} else {
			for k := range s.dirty {
				if dm, ok := any(s.data[k]).(iDirtyModel); ok {
					dm.CleanDirty()
				}
			}
		}
	}
	//即使data已被清空,也必须复位脏标记,否则BuildBson会重复吐出更新
	s.dirtyAll = false
	clear(s.dirty)
	clear(s.removed)
}

// updateDirty 更新藏标记
func (s *DMap[K, V]) updateDirty(tk any) {
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

// updateRemoved 标记key被删除(生成$unset)
func (s *DMap[K, V]) updateRemoved(k K) {
	if s.dirtyAll || s.removed[k] {
		return
	}
	s.removed[k] = true
	if s.dirtyParent != nil {
		s.dirtyParent.Invoke(s.inParentDirtyIdx)
	}
}

// updateDirtyAll 更新整个对象为脏
func (s *DMap[K, V]) updateDirtyAll() {
	if s.dirtyAll {
		return
	}
	s.dirtyAll = true
	if s.dirtyParent != nil {
		s.dirtyParent.Invoke(s.inParentDirtyIdx)
	}
}

// String toString
func (s *DMap[K, V]) String() string {
	return fmt.Sprintf("%v", s.data)
}

// MarshalJSON json序列化
func (s *DMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

// UnmarshalJSON json反序列化
func (s *DMap[K, V]) UnmarshalJSON(data []byte) error {
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
func (s *DMap[K, V]) MarshalBSON() ([]byte, error) {
	_, data, err := bson.MarshalValue(s.data)
	return data, err
}

// UnmarshalBSON bson反序列化
func (s *DMap[K, V]) UnmarshalBSON(data []byte) error {
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

// BuildBson bson的增量更新
func (s *DMap[K, V]) BuildBson(m bson.M, preKey string) {
	if len(s.dirty) == 0 && len(s.removed) == 0 && !s.dirtyAll {
		return
	}
	if s.dirtyAll {
		AddSetDirtyM(m, preKey, s)
		return
	}
	for k := range s.dirty {
		AddSetDirtyM(m, MakeBsonKey(fmt.Sprintf("%v", k), preKey), s.data[k])
	}
	for k := range s.removed {
		AddUnsetDirtyM(m, MakeBsonKey(fmt.Sprintf("%v", k), preKey))
	}
}

// ToMap to map
func (s *DMap[K, V]) ToMap() map[K]V {
	if s == nil || len(s.data) == 0 {
		return nil
	}
	ret := make(map[K]V, len(s.data))
	maps.Copy(ret, s.data)
	return ret
}
