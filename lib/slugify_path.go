package lib

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	unwantedCharsRegex = regexp.MustCompile(`[^\p{Hangul}a-zA-Z0-9\s]+`)
	whitespaceRegex    = regexp.MustCompile(`\s+`)
)

func SlugifyPath(filePath string) string {
	// 파일 경로에서 확장자를 제외한 순수 파일명을 가져옵니다
	base := filepath.Base(filePath)
	title := strings.TrimSuffix(base, filepath.Ext(base))

	// 허용되지 않는 모든 특수문자를 제거합니다
	safeTitle := unwantedCharsRegex.ReplaceAllString(title, "")

	// 하나 이상의 연속된 공백을 하이픈 1개로 변경합니다
	pathWithDashes := whitespaceRegex.ReplaceAllString(safeTitle, "-")

	// URL 일관성을 위해 모두 소문자로 변경합니다
	return strings.ToLower(pathWithDashes)
}
