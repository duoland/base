package strings

import "testing"

func TestInt64ToBase32(t *testing.T) {
	var src int64 = 1892882925901320192
	dst := Int64ToBase32(src)
	t.Log(dst)
}
