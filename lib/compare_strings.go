package lib

import "unicode"

const (
	typeNumber = iota
	typeKorean
	typeEnglish
	typeOther
)

func getCharType(r rune) int {
	if unicode.IsDigit(r) {
		return typeNumber
	}
	if r >= '가' && r <= '힣' {
		return typeKorean
	}
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return typeEnglish
	}
	return typeOther
}

// CompareStrings는 두 문자열을 '숫자 -> 한글 -> 영어' 순서 규칙에 따라 비교합니다.
// a가 b보다 앞에 와야 하면 true를 반환합니다.
func CompareStrings(a, b string) bool {
	runesA := []rune(a)
	runesB := []rune(b)

	lenA := len(runesA)
	lenB := len(runesB)

	minLen := lenA
	if lenB < minLen {
		minLen = lenB
	}

	for i := 0; i < minLen; i++ {
		runeA, runeB := runesA[i], runesB[i]
		typeA := getCharType(runeA)
		typeB := getCharType(runeB)

		if typeA != typeB {
			return typeA < typeB
		}
		if runeA != runeB {
			return runeA < runeB
		}
	}

	return lenA < lenB
}
