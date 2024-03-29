package common

// GenSpaces generate a space string of specified length
func GenSpaces(l int) string {
	return GenStr(l, " ")
}

// GenStr generate a string of the specified length, the string is composed of the given characters
func GenStr(l int, s string) string {
	var ss string
	for i := 0; i < l; i++ {
		ss += s
	}
	return ss
}