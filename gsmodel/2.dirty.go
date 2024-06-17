package gsmodel

type DirtyParentFunc func(dirtyIdx any)

func (f DirtyParentFunc) Invoke(dirtyIdx any) {
	if f != nil {
		f(dirtyIdx)
	}
}

// DirtyModel ------------------dirtyModel脏标记
type DirtyModel struct {
	dirty            uint64
	inParentDirtyIdx any
	dirtyParent      DirtyParentFunc
}

// SetParent 设置父节点
func (s *DirtyModel) SetParent(idx any, dirtyParentFunc DirtyParentFunc) {
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
func (s *DirtyModel) IsDirty() bool {
	return s.dirty > 0
}

// CleanDirty 清除脏标记
func (s *DirtyModel) CleanDirty() {
	if s == nil {
		return
	}
	d := s.dirty
	if d == 0 {
		return
	} else {
		s.dirty = 0
	}
}

// UpdateDirty 脏标记更新
func (s *DirtyModel) UpdateDirty(tn any) {
	n := uint64(tn.(int))
	if s.dirty&n == n {
		return
	}
	s.dirty |= n
	if s.dirtyParent != nil {
		s.dirtyParent.Invoke(s.inParentDirtyIdx)
	}
}

// GetDirty 获取脏值
func (s *DirtyModel) GetDirty() uint64 {
	return s.dirty
}
