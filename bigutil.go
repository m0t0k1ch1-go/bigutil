package bigutil

import (
	"strings"
)

func is0x(s string) bool {
	return s == "0x"
}

func has0xPrefix(s string) bool {
	return strings.HasPrefix(s, "0x")
}
