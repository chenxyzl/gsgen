package gsmodel

import "encoding/json"

// Clone 通过json序列化/反序列化深拷贝一个model,返回同类型的新指针。
// 生成代码统一调用此方法,避免每个类型重复生成相同的Clone逻辑。
// 注意:依赖类型自身生成的Marshal/UnmarshalJSON,不参与拷贝的字段(如脏标记)不会被复制。
func Clone[T any](s *T) (*T, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	ret := new(T)
	if err = json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
