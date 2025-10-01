package lib

import "unicode"

func CompareKoreanEnglish(s1, s2 string) bool {
	isKorean1 := isKorean([]rune(s1)[0])
	isKorean2 := isKorean([]rune(s2)[0])
	if isKorean1 != isKorean2 {
		return isKorean1 // 한글이 true이므로 앞으로 온다
	}
	return s1 < s2
}

func isKorean(r rune) bool { return unicode.Is(unicode.Hangul, r) }
