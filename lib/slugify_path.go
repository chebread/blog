package lib

import (
	"path/filepath"
	"strings"
)

// 웹 개발에서는 "A Go Programming Language"와 같은 제목을
// URL 친화적인 A-Go-Programming-Language 형태로 만든 것을 slug라고 부른다.
func SlugifyPath(filePath string) string {

	base := filepath.Base(filePath)
	title := strings.TrimSuffix(base, filepath.Ext(base))
	pathWithDashes := strings.ReplaceAll(title, " ", "-")

	return pathWithDashes
}
