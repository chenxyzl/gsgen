package gsmodel

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

// AMap ----------------------------------AMap-------------------------------------
// AMap map的包装
// @K key的类型
// @V value的类型
type AMap[K int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string, V any] struct {
	data map[K]V `bson:"map"`
}

func NewAMap[K int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string, V any]() *AMap[K, V] {
	ret := &AMap[K, V]{}
	ret.init()
	return ret
}
func (s *AMap[K, V]) init() {
	s.data = make(map[K]V)
}

// Len 长度
func (s *AMap[K, V]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

// Clean 重置清空list
func (s *AMap[K, V]) Clean() {
	if s == nil {
		return
	}
	s.data = make(map[K]V)
}

// Get 设置值
func (s *AMap[K, V]) Get(k K) V {
	if s == nil {
		panic("map is nil")
	}
	return s.data[k]
}

// Set 设置新值
func (s *AMap[K, V]) Set(k K, v V) {
	if s == nil {
		panic("map is nil")
	}
	s.data[k] = v
}

// Remove 删除
func (s *AMap[K, V]) Remove(k K) {
	if s == nil {
		panic("map is nil")
	}
	if _, ok := s.data[k]; !ok {
		return
	}
	delete(s.data, k)
}

// Range 遍历
func (s *AMap[K, V]) Range(f func(K, V) bool) {
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

// String toString
func (s *AMap[K, V]) String() string {
	return fmt.Sprintf("%v", s.data)
}

// MarshalJSON json序列化
func (s *AMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

// UnmarshalJSON json反序列化
func (s *AMap[K, V]) UnmarshalJSON(data []byte) error {
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
func (s *AMap[K, V]) MarshalBSON() ([]byte, error) {
	r, r1, r2 := bson.MarshalValue(s.data)
	_ = r
	return r1, r2
}

// UnmarshalBSON bson反序列化
func (s *AMap[K, V]) UnmarshalBSON(data []byte) error {
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

// ToMap to map
func (s *AMap[K, V]) ToMap() map[K]V {
	if len(s.data) == 0 {
		return nil
	}
	var ret = make(map[K]V)
	for k, v := range s.data {
		ret[k] = v
	}
	return ret
}
