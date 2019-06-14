package strings

func SliceContains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}


// ReverseRunes reverse a string runes and return the result
func ReverseRunes(str string) (output string) {
	if str == "" {
		return
	}

	runes:=[]rune(str)

	buf := make([]rune, 0, len(runes))
	for i := len(runes) - 1; i >= 0; i-- {
		buf = append(buf, runes[i])
	}
	output = string(buf)
	return
}
