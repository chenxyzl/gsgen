// Code generated by https://github.com/chenxyzl/gsgen; DO NOT EDIT.
// gen_tools version: 1.1.8
// generate time: 2024-07-03 16:03:35
package getter

import (
	"encoding/json"
	"fmt"
	"github.com/chenxyzl/gsgen/gsmodel"
)

func (s *TestA) GetCc() *gsmodel.AList[int] {
	return s.cc
}
func (s *TestA) SetCc(v *gsmodel.AList[int]) {
	s.cc = v
}
func (s *TestA) GetDd() *gsmodel.AMap[string, *TestA] {
	return s.dd
}
func (s *TestA) SetDd(v *gsmodel.AMap[string, *TestA]) {
	s.dd = v
}
func (s *TestA) String() string {
	doc := struct {
		Cc *gsmodel.AList[int]           `bson:"cc"`
		Dd *gsmodel.AMap[string, *TestA] `bson:"dd"`
	}{s.cc, s.dd}
	return fmt.Sprintf("%v", &doc)
}
func (s *TestA) MarshalJSON() ([]byte, error) {
	doc := struct {
		Cc *gsmodel.AList[int]           `bson:"cc"`
		Dd *gsmodel.AMap[string, *TestA] `bson:"dd"`
	}{s.cc, s.dd}
	return json.Marshal(doc)
}
func (s *TestA) UnmarshalJSON(data []byte) error {
	doc := struct {
		Cc *gsmodel.AList[int]           `bson:"cc"`
		Dd *gsmodel.AMap[string, *TestA] `bson:"dd"`
	}{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetCc(doc.Cc)
	s.SetDd(doc.Dd)
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
