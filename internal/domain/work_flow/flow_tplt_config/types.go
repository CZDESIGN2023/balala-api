package config

type Role struct {
	Id  int64  `json:"id"`
	Key string `json:"key"`
	Val string `json:"val"`
}

type Status struct {
	Id  string `json:"id"`
	Key string `json:"key"`
	Val string `json:"val"`
}
