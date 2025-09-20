package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func listFiles(dirPath string) []string {
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
					files = append(files, v.Name())
					fmt.Println(v.Name())
				}
			}
		}

		return files
	}

	f(dirPath) // 일단 처음 실행

	return files
}

func main() {
	// 모든 .md(x) 파일을 가져오기
	dirPath := "./content"
	files := listFiles(dirPath)
	fmt.Println(files)

	// 파일의 프론트 메터 부분만 읽기

	// 프론트 매터 가져오

	// published 분류하기

	// category 분류하기
}
