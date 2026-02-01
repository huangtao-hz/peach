package utils

import (
	"fmt"
	"regexp"
)

// Match 匹配字符串
func Match(pattern string, s string) (matched bool) {
	matched, err := regexp.MatchString(pattern, s)
	CheckFatal(err)
	return
}

// FullMatch 全匹配字符串
func FullMatch(pattern string, s string) bool {
	return Match(fmt.Sprintf("^%s$", pattern), s)
}

// Extract 提取匹配的字符串
func Extract(pattern, s string) string {
	r := regexp.MustCompile(pattern)
	return r.FindString(s)
}

// FindAll 提取所有匹配记录
func FindAll(pattern, s string) []string {
	r := regexp.MustCompile(pattern)
	return r.FindAllString(s, -1)
}
