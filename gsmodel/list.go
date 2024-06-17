package gsmodel

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

// AList ----------------------AList-----------------------
type AList[T any] struct {
	data []T `bson:"data"`
}

func NewAList[T any]() *AList[T] {
	ret := &AList[T]{}
	ret.init()
	return ret
}
func (s *AList[T]) init() {
	s.data = make([]T, 0)
}

// Len 长度
func (s *AList[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

// Clean 重置清空list
func (s *AList[T]) Clean() {
	if s == nil {
		return
	}
	s.data = make([]T, 0)
}

// Get 设置值
func (s *AList[T]) Get(idx int) T {
	if s == nil {
		panic("AList is nil")
	}
	l := s.Len()
	if idx >= l {
		panic(fmt.Sprintf("AList get idx out of range, len:%d|idx:%d", l, idx))
	}
	return s.data[idx]
}

// Set 设置新值
func (s *AList[T]) Set(idx uint64, v T) {
	if s == nil {
		panic("data is nil")
	}
	l := uint64(s.Len())
	if idx >= l {
		panic(fmt.Sprintf("AList set idx out of range, len:%d|idx:%d", l, idx))
	}
	s.data[idx] = v
}

// Append 追加
func (s *AList[T]) Append(vs ...T) {
	if s == nil {
		panic("data is nil")
	}
	for _, v := range vs {
		s.data = append(s.data, v)
	}
}

// Remove 删除
func (s *AList[T]) Remove(idx int) {
	if s == nil {
		panic("data is nil")
	}
	l := s.Len()
	if idx >= l {
		panic(fmt.Sprintf("AList remove idx out of range, len:%d|idx:%d", l, idx))
	}
	s.data = append(s.data[0:idx], s.data[idx+1:]...)
}

// Range 遍历
func (s *AList[T]) Range(f func(idx int, v T) bool) {
	if s == nil {
		panic("AList is nil")
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

// String toString
func (s *AList[T]) String() string {
	return fmt.Sprintf("%v", s.data)
}

// MarshalJSON json序列化
func (s *AList[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

// UnmarshalJSON json反序列化
func (s *AList[T]) UnmarshalJSON(data []byte) error {
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
func (s *AList[T]) MarshalBSON() ([]byte, error) {
	r, r1, r2 := bson.MarshalValue(s.data)
	_ = r
	return r1, r2
}

// UnmarshalBSON bson反序列化
func (s *AList[T]) UnmarshalBSON(data []byte) error {
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
