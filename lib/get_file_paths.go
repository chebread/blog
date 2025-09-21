package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFilePaths(dirPath string) []string {
	var files []string

	var f func(string) []string // 클로저 함수 타입 선언

	f = func(dirPath string) []string {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			fmt.Println(err)
		}

		for _, v := range entries {
			if v.IsDir() {
				// 파일이 아니라 디렉토리임.
				subDirPath := dirPath + "/" + v.Name()
				f(subDirPath)
			} else {
				if filepath.Ext(v.Name()) == ".md" || filepath.Ext(v.Name()) == ".mdx" {
					files = append(files, dirPath+"/"+v.Name())
				}
			}
		}

		return files
	}

	f(dirPath) // 일단 처음 실행

	return files
}
