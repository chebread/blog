package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

// 디렉터리가 존재하면 내용만 모두 비우고, 존재하지 않으면 해당 경로에 새 디렉터리를 생성함
func InitDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		fmt.Printf("정보: '%s' 디렉터리가 존재하지 않습니다. 새로 생성합니다.\n", dirPath)
		return os.MkdirAll(dirPath, 0755)
	}
	if err != nil {
		return err
	}

	fmt.Printf("성공: '%s' 디렉터리가 존재합니다. 내용물을 삭제합니다.\n", dirPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(dirPath, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
