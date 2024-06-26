// Code generated by https://github.com/chenxyzl/gsgen; DO NOT EDIT.
// gen_tools version: 1.1.5
// generate time: 2024-06-26 15:41:13
// test head annotations 1
// test head annotations 2
package nest

import (
	"encoding/json"
	"fmt"
	"github.com/chenxyzl/gsgen/gsmodel"
)

func (s *TestA) GetId() uint64 {
	return s.id
}
func (s *TestA) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(1 << 0)
}
func (s *TestA) GetAaa() int64 {
	return s.aaa
}
func (s *TestA) SetAaa(v int64) {
	s.aaa = v
	s.UpdateDirty(1 << 1)
}
func (s *TestA) GetBbb() int32 {
	return s.bbb
}
func (s *TestA) SetBbb(v int32) {
	s.bbb = v
	s.UpdateDirty(1 << 2)
}
func (s *TestA) GetCcc() string {
	return s.ccc
}
func (s *TestA) SetCcc(v string) {
	s.ccc = v
	s.UpdateDirty(1 << 3)
}
func (s *TestA) String() string {
	doc := struct {
		Id  uint64 `bson:"_id"`
		Aaa int64  `bson:"aaa"`
		Bbb int32  `bson:"bbb"`
		Ccc string `bson:"ccc"`
	}{s.id, s.aaa, s.bbb, s.ccc}
	return fmt.Sprintf("%v", &doc)
}
func (s *TestA) MarshalJSON() ([]byte, error) {
	doc := struct {
		Id  uint64 `bson:"_id"`
		Aaa int64  `bson:"aaa"`
		Bbb int32  `bson:"bbb"`
		Ccc string `bson:"ccc"`
	}{s.id, s.aaa, s.bbb, s.ccc}
	return json.Marshal(doc)
}
func (s *TestA) UnmarshalJSON(data []byte) error {
	doc := struct {
		Id  uint64 `bson:"_id"`
		Aaa int64  `bson:"aaa"`
		Bbb int32  `bson:"bbb"`
		Ccc string `bson:"ccc"`
	}{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetId(doc.Id)
	s.SetAaa(doc.Aaa)
	s.SetBbb(doc.Bbb)
	s.SetCcc(doc.Ccc)
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
func (s *TestB) GetId() uint64 {
	return s.id
}
func (s *TestB) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(1 << 0)
}
func (s *TestB) GetAa() string {
	return s.aa
}
func (s *TestB) SetAa(v string) {
	s.aa = v
	s.UpdateDirty(1 << 1)
}
func (s *TestB) GetBb() *TestA {
	return s.bb
}
func (s *TestB) SetBb(v *TestA) {
	if v != nil {
		v.SetParent(2, s.UpdateDirty)
	}
	s.bb = v
	s.UpdateDirty(1 << 2)
}
func (s *TestB) GetCc() *gsmodel.DList[*TestA] {
	return s.cc
}
func (s *TestB) SetCc(v *gsmodel.DList[*TestA]) {
	if v != nil {
		v.SetParent(3, s.UpdateDirty)
	}
	s.cc = v
	s.UpdateDirty(1 << 3)
}
func (s *TestB) GetDd() *gsmodel.DMap[string, *TestA] {
	return s.dd
}
func (s *TestB) SetDd(v *gsmodel.DMap[string, *TestA]) {
	if v != nil {
		v.SetParent(4, s.UpdateDirty)
	}
	s.dd = v
	s.UpdateDirty(1 << 4)
}
func (s *TestB) String() string {
	doc := struct {
		Id uint64                        `bson:"_id"`
		Aa string                        `bson:"aa"`
		Bb *TestA                        `bson:"bb"`
		Cc *gsmodel.DList[*TestA]        `bson:"cc"`
		Dd *gsmodel.DMap[string, *TestA] `bson:"dd"`
	}{s.id, s.aa, s.bb, s.cc, s.dd}
	return fmt.Sprintf("%v", &doc)
}
func (s *TestB) MarshalJSON() ([]byte, error) {
	doc := struct {
		Id uint64                        `bson:"_id"`
		Aa string                        `bson:"aa"`
		Bb *TestA                        `bson:"bb"`
		Cc *gsmodel.DList[*TestA]        `bson:"cc"`
		Dd *gsmodel.DMap[string, *TestA] `bson:"dd"`
	}{s.id, s.aa, s.bb, s.cc, s.dd}
	return json.Marshal(doc)
}
func (s *TestB) UnmarshalJSON(data []byte) error {
	doc := struct {
		Id uint64                        `bson:"_id"`
		Aa string                        `bson:"aa"`
		Bb *TestA                        `bson:"bb"`
		Cc *gsmodel.DList[*TestA]        `bson:"cc"`
		Dd *gsmodel.DMap[string, *TestA] `bson:"dd"`
	}{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetId(doc.Id)
	s.SetAa(doc.Aa)
	s.SetBb(doc.Bb)
	s.SetCc(doc.Cc)
	s.SetDd(doc.Dd)
	return nil
}
func (s *TestB) Clone() (*TestB, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	ret := TestB{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
func (s *TestC) GetId() uint64 {
	return s.id
}
func (s *TestC) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(1 << 0)
}
func (s *TestC) GetA() string {
	return s.a
}
func (s *TestC) SetA(v string) {
	s.a = v
	s.UpdateDirty(1 << 1)
}
func (s *TestC) GetB() *TestB {
	return s.b
}
func (s *TestC) SetB(v *TestB) {
	if v != nil {
		v.SetParent(2, s.UpdateDirty)
	}
	s.b = v
	s.UpdateDirty(1 << 2)
}
func (s *TestC) GetC() *gsmodel.DList[*TestB] {
	return s.c
}
func (s *TestC) SetC(v *gsmodel.DList[*TestB]) {
	if v != nil {
		v.SetParent(3, s.UpdateDirty)
	}
	s.c = v
	s.UpdateDirty(1 << 3)
}
func (s *TestC) GetD() *gsmodel.DMap[string, *TestB] {
	return s.d
}
func (s *TestC) SetD(v *gsmodel.DMap[string, *TestB]) {
	if v != nil {
		v.SetParent(4, s.UpdateDirty)
	}
	s.d = v
	s.UpdateDirty(1 << 4)
}
func (s *TestC) String() string {
	doc := struct {
		Id uint64                        `bson:"_id"`
		A  string                        `bson:"a"`
		B  *TestB                        `bson:"b"`
		C  *gsmodel.DList[*TestB]        `bson:"c"`
		D  *gsmodel.DMap[string, *TestB] `bson:"d"`
	}{s.id, s.a, s.b, s.c, s.d}
	return fmt.Sprintf("%v", &doc)
}
func (s *TestC) MarshalJSON() ([]byte, error) {
	doc := struct {
		Id uint64                        `bson:"_id"`
		A  string                        `bson:"a"`
		B  *TestB                        `bson:"b"`
		C  *gsmodel.DList[*TestB]        `bson:"c"`
		D  *gsmodel.DMap[string, *TestB] `bson:"d"`
	}{s.id, s.a, s.b, s.c, s.d}
	return json.Marshal(doc)
}
func (s *TestC) UnmarshalJSON(data []byte) error {
	doc := struct {
		Id uint64                        `bson:"_id"`
		A  string                        `bson:"a"`
		B  *TestB                        `bson:"b"`
		C  *gsmodel.DList[*TestB]        `bson:"c"`
		D  *gsmodel.DMap[string, *TestB] `bson:"d"`
	}{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetId(doc.Id)
	s.SetA(doc.A)
	s.SetB(doc.B)
	s.SetC(doc.C)
	s.SetD(doc.D)
	return nil
}
func (s *TestC) Clone() (*TestC, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	ret := TestC{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
