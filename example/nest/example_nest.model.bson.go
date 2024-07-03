// Code generated by https://github.com/chenxyzl/gsgen; DO NOT EDIT.
// gen_tools version: 1.1.8
// generate time: 2024-07-03 16:03:35
// test head annotations 1
// test head annotations 2
package nest

import (
	"github.com/chenxyzl/gsgen/gsmodel"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *TestA) MarshalBSON() ([]byte, error) {
	var doc = bson.M{"_id": s.id, "aaa": s.aaa, "bbb": s.bbb, "ccc": s.ccc}
	return bson.Marshal(doc)
}
func (s *TestA) UnmarshalBSON(data []byte) error {
	doc := struct {
		Id  uint64 `bson:"_id"`
		Aaa int64  `bson:"aaa"`
		Bbb int32  `bson:"bbb"`
		Ccc string `bson:"ccc"`
	}{}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetId(doc.Id)
	s.SetAaa(doc.Aaa)
	s.SetBbb(doc.Bbb)
	s.SetCcc(doc.Ccc)
	return nil
}
func (s *TestA) BuildBson(m bson.M, preKey string) {
	dirty := s.GetDirty()
	if dirty == 0 {
		return
	}
	if dirty&(1<<0) != 0 {
		gsmodel.AddSetDirtyM(m, gsmodel.MakeBsonKey("_id", preKey), s.id)
	}
	if dirty&(1<<1) != 0 {
		gsmodel.AddSetDirtyM(m, gsmodel.MakeBsonKey("aaa", preKey), s.aaa)
	}
	if dirty&(1<<2) != 0 {
		gsmodel.AddSetDirtyM(m, gsmodel.MakeBsonKey("bbb", preKey), s.bbb)
	}
	if dirty&(1<<3) != 0 {
		gsmodel.AddSetDirtyM(m, gsmodel.MakeBsonKey("ccc", preKey), s.ccc)
	}
	return
}
func (s *TestA) CleanDirty() {
	s.DirtyModel.CleanDirty()
}
func (s *TestB) MarshalBSON() ([]byte, error) {
	var doc = bson.M{"_id": s.id, "aa": s.aa, "bb": s.bb, "cc": s.cc, "dd": s.dd}
	return bson.Marshal(doc)
}
func (s *TestB) UnmarshalBSON(data []byte) error {
	doc := struct {
		Id uint64                        `bson:"_id"`
		Aa string                        `bson:"aa"`
		Bb *TestA                        `bson:"bb"`
		Cc *gsmodel.DList[*TestA]        `bson:"cc"`
		Dd *gsmodel.DMap[string, *TestA] `bson:"dd"`
	}{}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetId(doc.Id)
	s.SetAa(doc.Aa)
	s.SetBb(doc.Bb)
	s.SetCc(doc.Cc)
	s.SetDd(doc.Dd)
	return nil
}
func (s *TestB) BuildBson(m bson.M, preKey string) {
	dirty := s.GetDirty()
	if dirty == 0 {
		return
	}
	if dirty&(1<<0) != 0 {
		gsmodel.AddSetDirtyM(m, gsmodel.MakeBsonKey("_id", preKey), s.id)
	}
	if dirty&(1<<1) != 0 {
		gsmodel.AddSetDirtyM(m, gsmodel.MakeBsonKey("aa", preKey), s.aa)
	}
	if dirty&(1<<2) != 0 {
		if s.bb == nil {
			gsmodel.AddUnsetDirtyM(m, gsmodel.MakeBsonKey("bb", preKey))
		} else {
			s.bb.BuildBson(m, gsmodel.MakeBsonKey("bb", preKey))
		}
	}
	if dirty&(1<<3) != 0 {
		if s.cc == nil {
			gsmodel.AddUnsetDirtyM(m, gsmodel.MakeBsonKey("cc", preKey))
		} else {
			s.cc.BuildBson(m, gsmodel.MakeBsonKey("cc", preKey))
		}
	}
	if dirty&(1<<4) != 0 {
		if s.dd == nil {
			gsmodel.AddUnsetDirtyM(m, gsmodel.MakeBsonKey("dd", preKey))
		} else {
			s.dd.BuildBson(m, gsmodel.MakeBsonKey("dd", preKey))
		}
	}
	return
}
func (s *TestB) CleanDirty() {
	s.DirtyModel.CleanDirty()
	if s.bb != nil {
		s.bb.CleanDirty()
	}
	if s.cc != nil {
		s.cc.CleanDirty()
	}
	if s.dd != nil {
		s.dd.CleanDirty()
	}
}
func (s *TestC) MarshalBSON() ([]byte, error) {
	var doc = bson.M{"_id": s.id, "a": s.a, "b": s.b, "c": s.c, "d": s.d}
	return bson.Marshal(doc)
}
func (s *TestC) UnmarshalBSON(data []byte) error {
	doc := struct {
		Id uint64                        `bson:"_id"`
		A  string                        `bson:"a"`
		B  *TestB                        `bson:"b"`
		C  *gsmodel.DList[*TestB]        `bson:"c"`
		D  *gsmodel.DMap[string, *TestB] `bson:"d"`
	}{}
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	s.SetId(doc.Id)
	s.SetA(doc.A)
	s.SetB(doc.B)
	s.SetC(doc.C)
	s.SetD(doc.D)
	return nil
}
func (s *TestC) BuildBson(m bson.M, preKey string) {
	dirty := s.GetDirty()
	if dirty == 0 {
		return
	}
	if dirty&(1<<0) != 0 {
		gsmodel.AddSetDirtyM(m, gsmodel.MakeBsonKey("_id", preKey), s.id)
	}
	if dirty&(1<<1) != 0 {
		gsmodel.AddSetDirtyM(m, gsmodel.MakeBsonKey("a", preKey), s.a)
	}
	if dirty&(1<<2) != 0 {
		if s.b == nil {
			gsmodel.AddUnsetDirtyM(m, gsmodel.MakeBsonKey("b", preKey))
		} else {
			s.b.BuildBson(m, gsmodel.MakeBsonKey("b", preKey))
		}
	}
	if dirty&(1<<3) != 0 {
		if s.c == nil {
			gsmodel.AddUnsetDirtyM(m, gsmodel.MakeBsonKey("c", preKey))
		} else {
			s.c.BuildBson(m, gsmodel.MakeBsonKey("c", preKey))
		}
	}
	if dirty&(1<<4) != 0 {
		if s.d == nil {
			gsmodel.AddUnsetDirtyM(m, gsmodel.MakeBsonKey("d", preKey))
		} else {
			s.d.BuildBson(m, gsmodel.MakeBsonKey("d", preKey))
		}
	}
	return
}
func (s *TestC) CleanDirty() {
	s.DirtyModel.CleanDirty()
	if s.b != nil {
		s.b.CleanDirty()
	}
	if s.c != nil {
		s.c.CleanDirty()
	}
	if s.d != nil {
		s.d.CleanDirty()
	}
}
