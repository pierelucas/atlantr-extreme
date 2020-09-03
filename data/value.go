package data

import "strconv"

type value string

func (v *value) Int() int {
	s, _ := strconv.Atoi(string(*v))
	return s
}

func (v *value) String() string {
	return string(*v)
}

func (v *value) ToByte() []byte {
	return []byte(*v)
}
