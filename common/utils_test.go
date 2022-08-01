package common

import "testing"

func Test_interfaceToFloat(t *testing.T) {
	var want float64 = 1
	got := interfaceToFloat(int(1))
	if want != got {
		t.Error("wanted 1 as a float, but didn't get it")
	}
}
