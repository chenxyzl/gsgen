package model

import (
	"fmt"
	"gotest/model/mdata"
)

func (this *A) String() string {
	return fmt.Sprintf("x:%v,m:%v,n:%v", this.x, this.m, this.n)
}
func (this *A) GetX() int {
	return this.x
}
func (this *A) SetX(x int) {
	this.x = x
	this.UpdateDirty(0x01)
}
func (this *A) GetM() *mdata.MMap[int, int] {
	return &this.m
}
func (this *A) SetM(m *mdata.MMap[int, int]) {
	this.m = *m
	this.UpdateDirty(0x02)
	//
	//todo 非基本类型需要setSelfDirtyIdx
	//this.m.setSelfDirtyIdx(0x03, this.updateDirty)
}

func (this *A) GetY() *mdata.MList[int] {
	return &this.n
}
func (this *A) SetY(n *mdata.MList[int]) {
	this.n = *n
	this.UpdateDirty(0x03)
	//todo 非基本类型需要设置setSelfDirtyIdx
	//this.y.setSelfDirtyIdx(0x03, this.updateDirty)
}
