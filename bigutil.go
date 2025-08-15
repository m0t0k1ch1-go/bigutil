package bigutil

func is0x(s string) bool {
	return len(s) == 2 && s[0] == '0' && s[1] == 'x'
}

func has0xPrefix(s string) bool {
	return len(s) >= 2 && is0x(s[:2])
}
