package main

import (
	"gotest/model"
	"gotest/model/mdata"
)

func main() {
	a := &model.A{}
	a.SetX(1)
	println(a.GetX())
	a.SetY(mdata.NewList[int]())
	a.GetY().Append(2)
	println(a.GetY().Get(0))
	a.SetM(mdata.NewMMap[int, int]())
	b := a.GetM().Get(3)
	println(b)
	a.GetM().Set(3, 3)
	b = a.GetM().Get(3)
	println(b)
}
