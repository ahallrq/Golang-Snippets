package main

import "strings"

func Leftpad(s string, l int) string {
	padLen := l - len(s)
	if padLen < 0 { padLen = 0 }
	return strings.Repeat(" ", padLen) + s
}

func Rightpad(s string, l int) string {
	padLen := l - len(s)
	if padLen < 0 { padLen = 0 }
	return s + strings.Repeat(" ", padLen)
}
