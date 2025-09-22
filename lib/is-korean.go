package lib

import "unicode"

// 문자열의 첫 글자가 한글인지 판별하는 함수
func IsKorean(s string) bool {
	for _, r := range s {
		return unicode.Is(unicode.Hangul, r)
	}
	return false
}
