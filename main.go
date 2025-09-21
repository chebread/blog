package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getFilenames(dirPath string) []string {
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
					if v.Name() == "README.md" {
						// exclude file list
						// TODO: exclude 부분 개선해야 함.
						continue
					} else {
						files = append(files, v.Name())
					}
				}
			}
		}

		return files
	}

	f(dirPath) // 일단 처음 실행

	return files
}

type FrontMatter struct {
	Date      string   `yaml:"date"`
	Desc      string   `yaml:"desc"`
	Category  []string `yaml:"category"`
	Published bool     `yaml:"published"`
	Fixed     bool     `yaml:"fixed"`
}

func getFilePaths(dirPath string) []string {

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
					if v.Name() == "README.md" {
						// exclude file list
						// TODO: exclude 부분 개선해야 함.
						continue
					} else {
						files = append(files, dirPath+"/"+v.Name())
					}
				}
			}
		}

		return files
	}

	f(dirPath) // 일단 처음 실행

	return files
}

func main() {
	// 모든 .md(x) 파일명 가져오기
	var dirPath string = "./content"
	var filenames []string = getFilenames(dirPath)
	_ = filenames

	// 모든 .md(x) 파일 path 가져오기
	var filePaths []string = getFilePaths(dirPath)
	_ = filePaths

	// 프론트 매터 가져오기
	// published(날짜 별로 정렬), fixed(날짜 별로 정렬)에 file path 넣기
	var publishedFilePaths []string
	var fixedFilePaths []string
	_, _ = publishedFilePaths, fixedFilePaths // 슬라이스도 된다.

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("파일 열기 실패")
			continue
		}
		defer file.Close()

		buffer := make([]byte, 4096)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("%s - 파일 읽기 실패\n", filePath)
			fmt.Println()
			continue
		}
		content := buffer[:n]

		if !bytes.HasPrefix(content, []byte("---")) {
			fmt.Printf("%s - 프론트 매터를 찾을 수 없습니다\n", filePath)
			fmt.Println()
			continue
		}

		parts := bytes.SplitN(content, []byte("---"), 3)
		if len(parts) < 3 {
			fmt.Printf("%s - 프론트 매터의 끝을 찾을 수 없습니다\n", filePath)
			fmt.Println()
			continue
		}
		frontMatterBytes := parts[1]

		var fm FrontMatter
		if err := yaml.Unmarshal(frontMatterBytes, &fm); err != nil {
			fmt.Printf("%s - YAML 파싱 실패\n", filePath)
			fmt.Println()
			continue
		}

		// fmt.Printf("--- %s ---\n", filePath)
		// fmt.Printf("Date: %s\n", fm.Date)
		// fmt.Printf("Description: %s\n", fm.Desc)
		fmt.Printf("Category: %v\n", fm.Category)
		// fmt.Printf("Published: %v\n", fm.Published)
		// fmt.Printf("Fixed: %v\n", fm.Fixed)
		// fmt.Println()

		if fm.Published {
			if fm.Fixed {
				// fixed된 것은 어차피 published를 내포하고 있음으로, published에 넣지 않고 fixed에 넣는다
				fixedFilePaths = append(fixedFilePaths, filePath)
			} else {
				publishedFilePaths = append(publishedFilePaths, filePath)
			}
		}
		// TODO: 날짜 순서로 정렬하는 거는, 구조체, 맵 배우고 하자
	}

	fmt.Printf("fixed: %v\n", fixedFilePaths)
	fmt.Printf("published: %v\n", publishedFilePaths)

	// category 분류하기
}
