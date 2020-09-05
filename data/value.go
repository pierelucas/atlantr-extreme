package data

import "strconv"

type Value string

func (v *Value) Int() int {
	s, _ := strconv.Atoi(string(*v))
	return s
}

func (v *Value) String() string {
	return string(*v)
}

func (v *Value) ToByte() []byte {
	return []byte(*v)
}
