package main

import "testing"
import "strconv"

func TestSerializeDeserialize(t *testing.T) {

	for i := 1; i <= 20; i++ {
		var s string
		for j := 0; j < i; j++ {
			s = s + "x"
		}
		l := loginInfo{s}
		t.Log(l)
		serialized := serializeLoginInfo(&l)
		t.Log(serialized)
		result := deserializeLoginInfo(serialized)
		t.Log(result)

		if result.Username != l.Username {
			t.Error("The logininfo " + strconv.Itoa(i) + " did not serialize/deserialize correctly")
		}
	}
}
