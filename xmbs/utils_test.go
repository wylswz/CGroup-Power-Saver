package xmbs

import (
	"reflect"
	"testing"
)

func TestIn(t *testing.T) {
	enter := Enter([]int32{1, 2, 3, 4, 5}, []int32{4, 5, 6, 7, 8})
	if !reflect.DeepEqual([]int32{6, 7, 8}, enter) {
		t.Fail()
	}

	exit := Exit([]int32{1, 2, 3, 4, 5}, []int32{4, 5, 6, 7, 8})
	if !reflect.DeepEqual([]int32{1, 2, 3}, exit) {
		t.Fail()
	}
}

func doNothing(v interface{}) {

}

func doNothing2(v map[int]interface{}) {

}

func TestInterface(t *testing.T) {
	doNothing(1)

}

func TestMapKeys(t *testing.T) {
	m := map[int32]interface{}{
		1: "1",
		2: "2",
	}

	if !reflect.DeepEqual([]int32{1, 2}, Keys(m)) {
		t.Fail()
	}
}
