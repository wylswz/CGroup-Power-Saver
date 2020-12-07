package xmbs

import (
	"io/ioutil"
	"strconv"
)

// CheckErr Panic if err
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// MustReadFileAsBytes Read as byte array, or panic
func MustReadFileAsBytes(path string) []byte {
	content, err := ioutil.ReadFile(path)
	CheckErr(err)
	return content
}

// In if i is in il
func In(i int32, il []int32) bool {
	for _, v := range il {
		if i == v {
			return true
		}
	}
	return false
}

// Keys get keys of map
func Keys(m map[int32]interface{}) []int32 {
	ks := make([]int32, len(m))
	i := 0
	for k := range m {
		ks[i] = k
		i++
	}
	return ks
}

func Enter(old []int32, new []int32) []int32 {
	res := []int32{}
	for _, v := range new {
		if !In(v, old) {
			res = append(res, v)
		}
	}
	return res
}

func Exit(old []int32, new []int32) []int32 {
	res := []int32{}
	for _, v := range old {
		if !In(v, new) {
			res = append(res, v)
		}
	}
	return res
}

func ParseInt(s string) int64 {
	n, err := strconv.ParseInt(s, 0, 64)
	CheckErr(err)
	return n
}
