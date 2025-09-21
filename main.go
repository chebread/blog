package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

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
	// 모든 .md(x) 파일 path 가져오기
	var dirPath string = "./content"

	var filePaths []string = getFilePaths(dirPath)
	_ = filePaths

	// 프론트 매터 가져와서 published, fixed, category로 파일 경로 나누기
	var publishedFilePaths []string
	var fixedFilePaths []string
	var categoryFilePaths = make(map[string][]string)
	_, _, _ = publishedFilePaths, fixedFilePaths, categoryFilePaths // 슬라이스도 된다.

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
		// fmt.Printf("Category: %v\n", fm.Category)
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
			// TODO: 날짜 순서로 정렬하는 거는, 구조체, 맵 배우고 하자

			// category 뷴류: map의 슬라이스로 함.
			for _, c := range fm.Category {
				categoryFilePaths[c] = append(categoryFilePaths[c], filePath)
			}
		}
	}

	// 이거는 index.html의 블로그 리스트를 위한 슬라이스-맵 임.
	// fmt.Printf("fixed: %v\n", fixedFilePaths)
	// fmt.Printf("published: %v\n", publishedFilePaths)
	// fmt.Printf("categories: %v\n", categoryFilePaths)
	// TODO: 일단 카테고리 구현은 map 공부 이후 하자.

	// TODO: path에서 / 다 지우고 filename만 남기기.
	var fixedPostTitles [][]string
	var publishedPostTitles [][]string

	for _, filePath := range fixedFilePaths {
		// string slice
		// fmt.Println(v, strings.LastIndex(v, "/"), strings.LastIndex(v, "."))
		title := filePath[strings.LastIndex(filePath, "/")+1 : strings.LastIndex(filePath, ".")]
		fixedPostTitles = append(fixedPostTitles, []string{title, filePath})
	}
	for _, filePath := range publishedFilePaths {
		title := filePath[strings.LastIndex(filePath, "/")+1 : strings.LastIndex(filePath, ".")]
		publishedPostTitles = append(publishedPostTitles, []string{title, filePath})
	}

	// fmt.Println(fixedPostTitles, publishedPostTitles)

	// TODO: index.html의 post list 채우기
	// TODO: fixed / published / category를 html로서 각각 하나의 문자열 묶음처럼 저장해야 함.

	var htmlFixedPostList string = buildHtmlPostList(fixedPostTitles)
	var htmlPublishedPostList string = buildHtmlPostList(publishedPostTitles)

	fmt.Println(htmlFixedPostList)
	fmt.Println(htmlPublishedPostList)

	type PageData struct {
		Fixed     template.HTML
		Published template.HTML
	}

	data := PageData{
		Fixed:     template.HTML(htmlFixedPostList),
		Published: template.HTML(htmlPublishedPostList),
	}

	tmpl, err := template.ParseFiles("./layout/index.html")
	if err != nil {
		log.Fatalf("템플릿 파싱 실패: %v", err)
	}

	outputFile, err := os.Create("public/index.html")
	if err != nil {
		log.Fatalf("출력 파일 생성 실패: %v", err)
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, data); err != nil {
		log.Fatalf("템플릿 실행 실패: %v", err)
	}

	log.Println("성공: public/index.html 파일이 생성되었습니다.")

}

func buildHtmlPostList(srcSlice [][]string) string {
	var s []string
	s = append(s, "<ul>")
	for _, v := range srcSlice {
		title := v[0]
		path := v[1]
		html := fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", path, title)
		s = append(s, html)
	}
	s = append(s, "<ul>")
	html := strings.Join(s, "")
	return html
}
