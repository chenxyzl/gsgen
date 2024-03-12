package model

import (
	"gotest/model/mdata"
)

type A struct {
	//prop
	x int                  `bson:"x"`
	m mdata.MMap[int, int] `bson:"m"`
	n mdata.MList[int]     `bson:"m"`
	//------dirty------
	mdata.DirtyModel
}

type B struct {
	x int `bson:"x"`
	//------dirty------
	mdata.DirtyModel
}
