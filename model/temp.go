package model

import "encoding/json"

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
