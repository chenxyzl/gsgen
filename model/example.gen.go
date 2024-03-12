package model

func (s *TestA) GetId() uint64 {
	return s.id
}
func (s *TestA) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(0)
}
func (s *TestA) GetA() int64 {
	return s.a
}
func (s *TestA) SetA(v int64) {
	s.a = v
	s.UpdateDirty(1)
}
func (s *TestA) GetB() int32 {
	return s.b
}
func (s *TestA) SetB(v int32) {
	s.b = v
	s.UpdateDirty(2)
}
func (s *TestB) GetId() uint64 {
	return s.id
}
func (s *TestB) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(0)
}
func (s *TestB) GetM() string {
	return s.m
}
func (s *TestB) SetM(v string) {
	s.m = v
	s.UpdateDirty(1)
}
func (s *TestB) GetN() *TestA {
	return s.n
}
func (s *TestB) SetN(v *TestA) {
	s.n = v
	s.UpdateDirty(2)
	if v != nil {
		v.SetSelfDirtyIdx(2, s.UpdateDirty)
	}
}
func (s *TestC) GetId() uint64 {
	return s.id
}
func (s *TestC) SetId(v uint64) {
	s.id = v
	s.UpdateDirty(0)
}
func (s *TestC) GetM() string {
	return s.m
}
func (s *TestC) SetM(v string) {
	s.m = v
	s.UpdateDirty(1)
}
func (s *TestC) GetN() *TestB {
	return s.n
}
func (s *TestC) SetN(v *TestB) {
	s.n = v
	s.UpdateDirty(2)
	if v != nil {
		v.SetSelfDirtyIdx(2, s.UpdateDirty)
	}
}
