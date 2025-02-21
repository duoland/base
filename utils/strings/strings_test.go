package strings

import "testing"

func TestInt64ToBase32(t *testing.T) {
	var src int64 = 1892882925901320192
	dst := Int64ToBase32(src)
	t.Log("base32=", dst)
	src1, _ := Base32ToInt64(dst)
	if src1 != src {
		t.Fatal("base32 convert fail")
	}
	t.Log("int64=", src1)
}
