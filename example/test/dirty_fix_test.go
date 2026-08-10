package test

import (
	"testing"

	"github.com/chenxyzl/gsgen/example/nest"
	"github.com/chenxyzl/gsgen/gsmodel"
	"go.mongodb.org/mongo-driver/bson"
)

// TestDListBasicCleanDirty 验证基本类型元素的DList在Set后CleanDirty不会panic(修复#1)
func TestDListBasicCleanDirty(t *testing.T) {
	l := gsmodel.NewDList[int]()
	l.Append(1, 2, 3)
	l.Set(1, 20) //标记idx=1为脏
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("DList[int].CleanDirty 不应panic: %v", r)
		}
	}()
	l.CleanDirty()
	if l.IsDirty() {
		t.Fatal("CleanDirty后不应仍为脏")
	}
}

// TestDListCleanAfterEmpty 验证清空后脏标记也能被复位(修复#3)
func TestDListCleanAfterEmpty(t *testing.T) {
	l := gsmodel.NewDList[int]()
	l.Append(1, 2, 3)
	l.Clean() //清空并置dirtyAll
	if !l.IsDirty() {
		t.Fatal("Clean后应为脏")
	}
	l.CleanDirty()
	if l.IsDirty() {
		t.Fatal("清空后CleanDirty应能复位脏标记")
	}
	m := bson.M{}
	l.BuildBson(m, "k")
	if len(m) != 0 {
		t.Fatalf("CleanDirty后BuildBson不应再产生更新, got:%v", m)
	}
}

// TestDMapRemoveUnset 验证删除key生成$unset而非$set零值(修复#2)
func TestDMapRemoveUnset(t *testing.T) {
	mp := gsmodel.NewDMap[string, int]()
	mp.Set("a", 1)
	mp.Set("b", 2)
	mp.CleanDirty()

	mp.Remove("a")
	m := bson.M{}
	mp.BuildBson(m, "root")

	unset, ok := m["$unset"].(bson.M)
	if !ok {
		t.Fatalf("删除key应生成$unset, got:%v", m)
	}
	if _, ok := unset["root.a"]; !ok {
		t.Fatalf("$unset应包含root.a, got:%v", unset)
	}
	if _, ok := m["$set"]; ok {
		t.Fatalf("纯删除不应产生$set, got:%v", m)
	}
}

// TestDMapReAddCancelsRemove 验证删除后重新赋值会取消$unset,改为$set(修复#2)
func TestDMapReAddCancelsRemove(t *testing.T) {
	mp := gsmodel.NewDMap[string, int]()
	mp.Set("a", 1)
	mp.CleanDirty()

	mp.Remove("a")
	mp.Set("a", 9) //重新赋值

	m := bson.M{}
	mp.BuildBson(m, "root")

	if _, ok := m["$unset"]; ok {
		t.Fatalf("重新赋值后不应有$unset, got:%v", m)
	}
	set, ok := m["$set"].(bson.M)
	if !ok || set["root.a"] != 9 {
		t.Fatalf("重新赋值后应有$set root.a=9, got:%v", m)
	}
}

// TestDMapCleanAfterEmpty 验证map清空后脏标记能复位(修复#3)
func TestDMapCleanAfterEmpty(t *testing.T) {
	mp := gsmodel.NewDMap[string, int]()
	mp.Set("a", 1)
	mp.Clean()
	mp.CleanDirty()
	if mp.IsDirty() {
		t.Fatal("清空后CleanDirty应能复位脏标记")
	}
}

// TestDListIsDirtyAfterRemove 验证Remove后IsDirty返回true(dirtyAll场景)
func TestDListIsDirtyAfterRemove(t *testing.T) {
	l := gsmodel.NewDList[int]()
	l.Append(1, 2, 3)
	l.CleanDirty()
	if l.IsDirty() {
		t.Fatal("CleanDirty后不应为脏")
	}
	l.Remove(0) //Remove置dirtyAll,dirty map为空
	if !l.IsDirty() {
		t.Fatal("Remove后IsDirty应为true(此前因漏判dirtyAll而返回false)")
	}
}

// TestNilContainerIsDirty 验证nil容器IsDirty不panic且返回false
func TestNilContainerIsDirty(t *testing.T) {
	var l *gsmodel.DList[int]
	var mp *gsmodel.DMap[string, int]
	if l.IsDirty() || mp.IsDirty() {
		t.Fatal("nil容器不应为脏")
	}
}

// TestGeneratedCloneViaHelper 验证生成的Clone(委托给gsmodel.Clone)是深拷贝且相互独立
func TestGeneratedCloneViaHelper(t *testing.T) {
	c := getTestC(123)
	clone, err := c.Clone()
	if err != nil {
		t.Fatalf("Clone失败: %v", err)
	}
	if clone == c {
		t.Fatal("Clone应返回新指针")
	}
	if clone.GetId() != c.GetId() || clone.GetA() != c.GetA() {
		t.Fatalf("Clone字段应一致: id=%d/%d a=%s/%s", clone.GetId(), c.GetId(), clone.GetA(), c.GetA())
	}
	//修改副本不应影响原对象
	clone.SetA("changed")
	if c.GetA() == "changed" {
		t.Fatal("修改Clone不应影响原对象(深拷贝)")
	}
}

// TestCloneHelperDirect 直接验证gsmodel.Clone泛型helper
func TestCloneHelperDirect(t *testing.T) {
	a := &nest.TestA{}
	a.SetId(7)
	a.SetCcc("x")
	b, err := gsmodel.Clone(a)
	if err != nil {
		t.Fatalf("gsmodel.Clone失败: %v", err)
	}
	if b == a || b.GetId() != 7 || b.GetCcc() != "x" {
		t.Fatalf("gsmodel.Clone结果错误: %+v", b)
	}
}
