package test1

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TA[T any] struct {
	x string `json:"x,omitempty"`
}

func (s TA[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		X string `json:"x,omitempty"`
	}{X: s.x})
}

func (s *TA[T]) UnmarshalJSON(data []byte) error {
	v := struct {
		X string `json:"x,omitempty"`
	}{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.x = v.X
	return nil
}
func (this *TA[T]) F() {
	fmt.Println("xx")
}

func TestF(t *testing.T) {
	var a TA[int]
	a.F()
	var b = new(TA[int])
	b.F()

	a.x = "1111"
	v, err := json.Marshal(a)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(1)
		fmt.Println(string(v))
		fmt.Println(2)
		fmt.Println(a)
		fmt.Println(3)
	}
	var c TA[int]
	err = json.Unmarshal(v, &c)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(c)
		fmt.Println(4)
	}
}
func TestF1(t *testing.T) {
	var a int

	a = 1111
	v, err := json.Marshal(a)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(1)
		fmt.Println(string(v))
	}
	var c int
	err = json.Unmarshal(v, &c)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(c)
		fmt.Println(4)
	}
}
