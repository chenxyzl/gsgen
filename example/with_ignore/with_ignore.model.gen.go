// Code generated by https://github.com/chenxyzl/gsgen; DO NOT EDIT.
// gen_tools version: 1.1.8
// generate time: 2024-07-03 16:03:36
package with_ignore

import (
	"encoding/json"
	"fmt"
	"github.com/chenxyzl/gsgen/example/common"
)

func (s *TestA) GetIg() *common.Common {
	return s.ig
}
func (s *TestA) SetIg(v *common.Common) {
	if v != nil {
		v.SetParent(0, s.UpdateDirty)
	}
	s.ig = v
	s.UpdateDirty(1 << 0)
}
func (s *TestA) String() string {
	doc := struct {
		Ig *common.Common `bson:"ig"`
	}{s.ig}
	return fmt.Sprintf("%v", &doc)
}
func (s *TestA) MarshalJSON() ([]byte, error) {
	doc := struct {
		Ig *common.Common `bson:"ig"`
	}{s.ig}
	return json.Marshal(doc)
}
func (s *TestA) UnmarshalJSON(data []byte) error {
	doc := struct {
		Ig *common.Common `bson:"ig"`
	}{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetIg(doc.Ig)
	return nil
}
func (s *TestA) Clone() (*TestA, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	ret := TestA{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
