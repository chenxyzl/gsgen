package gsmodel

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

// DList ----------------------DList-----------------------
type DList[T any] struct {
	data []T `bson:"data"`
	//
	dirty            map[uint64]bool
	dirtyAll         bool
	inParentDirtyIdx any
	dirtyParent      DirtyParentFunc
}

func NewDList[T any]() *DList[T] {
	ret := &DList[T]{}
	ret.init()
	return ret
}
func (s *DList[T]) init() {
	s.data = make([]T, 0)
	s.dirty = make(map[uint64]bool)
}

// Len 长度
func (s *DList[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

// Clean 重置清空list
func (s *DList[T]) Clean() {
	if s == nil || len(s.data) == 0 {
		return
	}
	s.data = make([]T, 0)
	s.updateDirtyAll()
}

// Get 设置值
func (s *DList[T]) Get(idx int) T {
	if s == nil {
		panic("DList is nil")
	}
	l := s.Len()
	if idx >= l {
		panic(fmt.Sprintf("DList get idx out of range, len:%d|idx:%d", l, idx))
	}
	return s.data[idx]
}

// Set 设置新值
func (s *DList[T]) Set(idx uint64, v T) {
	if s == nil {
		panic("data is nil")
	}
	l := uint64(s.Len())
	if idx >= l {
		panic(fmt.Sprintf("DList set idx out of range, len:%d|idx:%d", l, idx))
	}
	//
	checkSetParent(v, idx, s.updateDirty)
	s.data[idx] = v
	s.updateDirty(idx)
}

// Append 追加
func (s *DList[T]) Append(vs ...T) {
	if s == nil {
		panic("data is nil")
	}
	for _, v := range vs {
		idx := uint64(s.Len())
		//
		checkSetParent(v, idx, s.updateDirty)
		s.data = append(s.data, v)
		s.updateDirty(idx)
	}
}

// Remove 删除 注:因为删除不太好处理list对应的bson的更新,所以这里用了DirtyAll
func (s *DList[T]) Remove(idx int) {
	if s == nil {
		panic("data is nil")
	}
	l := s.Len()
	if idx >= l {
		panic(fmt.Sprintf("DList remove idx out of range, len:%d|idx:%d", l, idx))
	}
	s.data = append(s.data[0:idx], s.data[idx+1:]...)
	s.updateDirtyAll()
}

// Range 遍历
func (s *DList[T]) Range(f func(idx int, v T) bool) {
	if s == nil {
		panic("DList is nil")
	}
	if f == nil {
		return
	}
	for idx, v := range s.data {
		if _continue := f(idx, v); !_continue {
			break
		}
	}
}

// SetParent 设置父节点
func (s *DList[T]) SetParent(idx any, dirtyParentFunc DirtyParentFunc) {
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
func (s *DList[T]) IsDirty() bool {
	return len(s.dirty) > 0
}

// CleanDirty 清楚脏标记
func (s *DList[T]) CleanDirty() {
	if s == nil || len(s.data) == 0 {
		return
	}
	if s.dirtyAll {
		var v T
		if _, ok := (any(v)).(IDirtyModel); ok {
			s.Range(func(idx int, v T) bool {
				(any(v)).(IDirtyModel).CleanDirty()
				return true
			})
		}
	} else {
		l := s.Len()
		for idx, dirty := range s.dirty {
			if dirty && int(idx) < l {
				(any(s.Get(int(idx)))).(IDirtyModel).CleanDirty()
			}
		}
	}
	s.dirtyAll = false
	clear(s.dirty)
}

// updateDirty 标记脏
func (s *DList[T]) updateDirty(a any) {
	n := a.(uint64)
	//如果已经allDirty了就不用管了
	if s.dirtyAll || s.dirty[n] {
		return
	}
	s.dirty[n] = true
	if s.dirtyParent != nil {
		s.dirtyParent.Invoke(s.inParentDirtyIdx)
	}
}

// updateDirtyAll 标记所有都为脏
func (s *DList[T]) updateDirtyAll() {
	if s.dirtyAll {
		return
	}
	s.dirtyAll = true
	if s.dirtyParent != nil {
		s.dirtyParent.Invoke(s.inParentDirtyIdx)
	}
}

// String toString
func (s *DList[T]) String() string {
	return fmt.Sprintf("%v", s.data)
}

// MarshalJSON json序列化
func (s *DList[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

// UnmarshalJSON json反序列化
func (s *DList[T]) UnmarshalJSON(data []byte) error {
	var list []T
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	s.init()
	for _, v := range list {
		s.Append(v)
	}
	return nil
}

// MarshalBSON bson序列化
func (s *DList[T]) MarshalBSON() ([]byte, error) {
	r, r1, r2 := bson.MarshalValue(s.data)
	_ = r
	return r1, r2
}

// UnmarshalBSON bson反序列化
func (s *DList[T]) UnmarshalBSON(data []byte) error {
	var list []T
	if err := bson.UnmarshalValue(bson.TypeArray, data, &list); err != nil {
		return err
	}
	s.init()
	for _, v := range list {
		s.Append(v)
	}
	return nil
}

// BuildBson bson的增量更新
func (s *DList[T]) BuildBson(m bson.M, preKey string) {
	if len(s.dirty) == 0 && !s.dirtyAll {
		return
	}
	if s.dirtyAll {
		AddSetDirtyM(m, preKey, s)
	} else {
		for idx := range s.dirty {
			AddSetDirtyM(m, MakeBsonKey(fmt.Sprintf("%d", idx), preKey), s.data[idx])
		}
	}
	return
}
