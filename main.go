package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"

	"blog/lib"
)

func buildHtmlPostList(srcSlice [][]string) string {
	var s []string
	s = append(s, "<ul>")
	for _, v := range srcSlice {
		title := v[0]
		path := v[1]
		html := fmt.Sprintf("<li><a href=\"post/%s\">%s</a></li>", path, title)
		s = append(s, html)
	}
	s = append(s, "<ul>")
	html := strings.Join(s, "")
	return html
}

func filePathToURLPath(filePath string) string {
	base := filepath.Base(filePath)
	title := strings.TrimSuffix(base, filepath.Ext(base))
	pathWithDashes := strings.ReplaceAll(title, " ", "-")
	urlPath := url.PathEscape(pathWithDashes)

	return fmt.Sprintf("%s.html", urlPath)
}

func main() {

	// 모든 .md(x) 파일 path 가져오기
	var dirPath string = "./content"

	var filePaths []string = lib.GetFilePaths(dirPath)

	// 프론트 매터 가져와서 published, fixed, category로 파일 경로 나누기
	type FrontMatter struct {
		Date      string   `yaml:"date"`
		Desc      string   `yaml:"desc"`
		Category  []string `yaml:"category"`
		Published bool     `yaml:"published"`
		Fixed     bool     `yaml:"fixed"`
	}

	// 이거는 index.html의 블로그 리스트를 위한 슬라이스-맵 임.
	var publishedFilePaths []string
	var fixedFilePaths []string
	var categoryFilePaths = make(map[string][]string)

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
		} else {
			fmt.Println(filePath)
		}
	}

	// path에서 filename만 남기기.

	// TODO: 일단 카테고리는 map 공부 이후 하자.

	var fixedPostInfos [][]string
	var publishedPostInfos [][]string

	for _, filePath := range fixedFilePaths {
		title := filePath[strings.LastIndex(filePath, "/")+1 : strings.LastIndex(filePath, ".")]
		path := filePathToURLPath(title)
		fixedPostInfos = append(fixedPostInfos, []string{title, path, filePath})
	}
	for _, filePath := range publishedFilePaths {
		title := filePath[strings.LastIndex(filePath, "/")+1 : strings.LastIndex(filePath, ".")]
		path := filePathToURLPath(title)
		publishedPostInfos = append(publishedPostInfos, []string{title, path, filePath})
	}

	// index.html의 post list 채우기
	var htmlFixedPostList string = buildHtmlPostList(fixedPostInfos)
	var htmlPublishedPostList string = buildHtmlPostList(publishedPostInfos)

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
		fmt.Printf("템플릿 파싱 실패\n")
	}

	outputFile, err := os.Create("public/index.html")
	if err != nil {
		fmt.Printf("출력 파일 생성 실패\n")
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, data); err != nil {
		fmt.Printf("템플릿 실행 실패\n")
	}
	log.Printf("성공: [index.html] 파일이 생성되었습니다.\n")

	// Post 생성하기
	var postDirPath string = "public/post"

	// public/post 디렉토리 초기화
	if err := os.RemoveAll(postDirPath); err != nil {
		fmt.Printf("public/post 디렉토리 삭제 실패\n")
	}

	// public/post 디렉토리 생성
	if err := os.MkdirAll(postDirPath, 0755); err != nil {
		fmt.Printf("public/post 디렉토리 생성 실패\n")
	}

	var postInfos [][]string = append(fixedPostInfos, publishedPostInfos...) // 참고로, 여기에는 category는 포함시키지 않아도 됨.

	md := goldmark.New(
		goldmark.WithExtensions(
			meta.New(),
		),
	)

	layoutFile := "layout/post.html"
	tmplPost, err := template.ParseFiles(layoutFile)
	if err != nil {
		log.Fatalf("템플릿 파싱 실패: %v", err)
	}

	type PostData struct {
		Title   string
		Date    string
		Content template.HTML
	}

	for _, postInfo := range postInfos {
		title := postInfo[0]
		filename := postInfo[1]
		filePath := postInfo[2]

		sourceBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("파일 읽기 실패\n")
			continue
		}

		var contentBuf bytes.Buffer
		context := parser.NewContext()
		if err := md.Convert(sourceBytes, &contentBuf, parser.WithContext(context)); err != nil {
			fmt.Printf("Markdown 변환 실패\n")
			continue
		}

		metaData := meta.Get(context)
		var fm FrontMatter
		if metaData != nil {
			metaBytes, _ := yaml.Marshal(metaData)
			if err := yaml.Unmarshal(metaBytes, &fm); err != nil {
				fmt.Printf("프론트 메터 파싱 실패\n")
				// TODO: Date를 또 따로 파싱 하지 말고, 위에서 yaml 할 때 한 번에 파싱 하자.
				continue
			}
		}

		page := PostData{
			Title:   title,
			Date:    fm.Date,
			Content: template.HTML(contentBuf.String()),
		}

		decodedFilename, err := url.PathUnescape(filename)
		if err != nil {
			// 디코딩 실패 시 에러 처리 (예: 잘못된 % 인코딩)
			fmt.Printf("파일명 디코딩 실패\n")
			continue
		}

		outputPath := filepath.Join("public", "post", decodedFilename)
		outputFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("출력 파일 생성 실패\n")
			continue
		}
		defer outputFile.Close()

		if err := tmplPost.Execute(outputFile, page); err != nil {
			fmt.Printf("템플릿 실행 실패\n")
			continue
		}

		log.Printf("성공: [%s] 파일이 생성되었습니다.", outputPath)
	}

	// CSS를 위한 layout/style 디렉토리 복사
	sourceStylesDir := "layout/styles"
	destStylesDir := "public/styles"
	if err := lib.CopyDir(sourceStylesDir, destStylesDir); err != nil {
		fmt.Printf("layout/style 디렉토리 복사 실패\n")
	}

	// favicon을 위한 favicons 디렉토리 복사
	sourceFaviconsDir := "layout/favicons"
	destFaviconsDir := "public/favicons"
	if err := lib.CopyDir(sourceFaviconsDir, destFaviconsDir); err != nil {
		fmt.Printf("layout/favicons 디렉토리 복사 실패\n")
	}

	// TODO: github action으로 git tag ... 이렇게 하면 go run 해서 public에 정적 파일 생성해서 그거를 github.io에 올리면 됨.
}
