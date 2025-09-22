package main

import (
	"blog/lib"

	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

// TODO: IO 공부하기

func main() {
	// *** post data 처리
	var contentDirPath string = "./content"
	var contentFilePaths []string = lib.GetFilePaths(contentDirPath)
	_ = contentFilePaths

	var postsData []map[string]any                              // published된 모든 post의 정보
	var postsDataByCategory = make(map[string][]map[string]any) // 카테고리로 분류한 published된 모든 post의 정보
	_ = postsData
	_ = postsDataByCategory

	type PostFrontMatter struct {
		Date      string   `yaml:"date"`
		Desc      string   `yaml:"desc"`
		Category  []string `yaml:"category"`
		Published bool     `yaml:"published"`
		Fixed     bool     `yaml:"fixed"`
	}

	for _, path := range contentFilePaths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("파일 열기 실패")
			continue
		}
		defer file.Close()

		buffer := make([]byte, 4096)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("%s - 파일 읽기 실패\n", path)
			fmt.Println()
			continue
		}
		content := buffer[:n]

		if !bytes.HasPrefix(content, []byte("---")) {
			fmt.Printf("%s - 프론트 매터를 찾을 수 없습니다\n", path)
			fmt.Println()
			continue
		}

		parts := bytes.SplitN(content, []byte("---"), 3)
		if len(parts) < 3 {
			fmt.Printf("%s - 프론트 매터의 끝을 찾을 수 없습니다\n", path)
			fmt.Println()
			continue
		}
		frontMatterBytes := parts[1]

		var fm PostFrontMatter
		if err := yaml.Unmarshal(frontMatterBytes, &fm); err != nil {
			fmt.Printf("%s - YAML 파싱 실패\n", path)
			fmt.Println()
			continue
		}

		if fm.Published {
			var title string = path[strings.LastIndex(path, "/")+1 : strings.LastIndex(path, ".")]
			var outputFilePath string = fmt.Sprintf("%s.html", lib.SlugifyPath(title)) // filepath.Join("public", "post", ...)  이거 않하고 그냥 따로 처리

			var description string
			if fm.Desc != "" {
				description = fm.Desc
			} else {
				bodyContent := bytes.TrimSpace(parts[2])

				plainText, err := lib.MarkdownToPlainText(bodyContent)
				if err != nil {
					log.Printf("Markdown 변환 실패: %v", err)
					description = ""
				} else {
					normalizedText := strings.Join(strings.Fields(plainText), " ")

					runes := []rune(normalizedText)
					if len(runes) > 100 {
						description = string(runes[:100]) + "..."
					} else {
						description = string(runes)
					}
				}
			}

			var data = map[string]any{
				"title":          title,
				"outputFilePath": outputFilePath, // dash만 되어 있고, 비인코딩되어 있어야 함.
				"sourceFilePath": path,
				"category":       fm.Category,
				"date":           fm.Date,
				"fixed":          fm.Fixed,
				"description":    description, // 만약 desc가 없으면 공백("")이 저장됨.
			}

			postsData = append(postsData, data)

			for _, category := range fm.Category {
				// map -> slice -> map
				postsDataByCategory[category] = append(postsDataByCategory[category], data)
			}
		}
	}

	// *** public 디렉토리 삭제 및 생성
	publicDir := "public"
	if err := lib.InitDir(publicDir); err != nil {
		log.Fatalf("public 디렉토리 초기화 실패: %v", err)
	}

	// public 위한 파일, 디렉토리 복사
	sourceStylesDir := "layout/styles"
	destStylesDir := "public/styles"
	if err := lib.CopyDir(sourceStylesDir, destStylesDir); err != nil {
		fmt.Printf("layout/style 디렉토리 복사 실패\n")
	}
	fmt.Printf("성공: layout/style 디렉토리 복사\n")

	sourceFaviconsDir := "layout/favicons"
	destFaviconsDir := "public/favicons"
	if err := lib.CopyDir(sourceFaviconsDir, destFaviconsDir); err != nil {
		fmt.Printf("layout/favicons 디렉토리 복사 실패\n")
	}
	fmt.Printf("성공: layout/favicons 디렉토리 복사\n")

	sourceAssetsDir := "content/assets"
	destAssetsDir := "public/assets"
	if err := lib.CopyDir(sourceAssetsDir, destAssetsDir); err != nil {
		fmt.Printf("content/assets 디렉토리 복사 실패\n")
	}
	fmt.Printf("성공: content/assets 디렉토리 복사\n")

	// gp-pages에서 기본적으로 제공하는 404 사용
	// source404File := "layout/404.html"
	// dest404File := "public/404.html"
	// if err := lib.CopyFile(source404File, dest404File); err != nil {
	// 	fmt.Printf("layout/404.html 파일 복사 실패\n")
	// }
	// fmt.Printf("성공: layout/404.html 파일 복사\n")

	noJekyllPath := filepath.Join("public", ".nojekyll")
	if err := os.WriteFile(noJekyllPath, []byte(""), 0644); err != nil {
		fmt.Printf(".nojekyll 파일 생성 실패\n")
	}
	fmt.Printf("성공: public/.nojekyll 파일 생성\n")

	// *** category 기반 post list 처리
	var postList []string

	// 정렬: 이름순(한글 -> 영어)
	var categories []string
	for category := range postsDataByCategory {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		keyI := categories[i]
		keyJ := categories[j]
		isKoreanI := lib.IsKorean(keyI)
		isKoreanJ := lib.IsKorean(keyJ)

		if isKoreanI != isKoreanJ {
			return isKoreanI
		}

		return keyI < keyJ
	})

	for _, category := range categories {
		postsData := postsDataByCategory[category]
		postList = append(postList, fmt.Sprintf("<h2>%s</h2>", category))
		postList = append(postList, "<ul>")

		// 정렬: fixed, date
		sort.Slice(postsData, func(i, j int) bool {
			postI := postsData[i]
			postJ := postsData[j]

			fixedI, _ := postI["fixed"].(bool)
			fixedJ, _ := postJ["fixed"].(bool)

			if fixedI != fixedJ {
				return fixedI
			}

			dateI, _ := postI["date"].(string)
			dateJ, _ := postJ["date"].(string)
			return dateI > dateJ
		})

		for _, data := range postsData {
			isFixed, _ := data["fixed"].(bool)
			outputFilePath, _ := data["outputFilePath"].(string)
			permalink := filepath.Join("post", outputFilePath) // 블로그 링크

			if isFixed {
				fixedTemplate := `<li>
					<article>
						<h3 class="post-title"><a href="%s">[Fixed] %s</a></h3>
						<p class="post-date"><time datetime="%s">%s</time></p>
						<p class="post-description">%s</p>
						<div class="post-category">%s</div>
					</article>
				</li>`
				postList = append(postList, fmt.Sprintf(fixedTemplate, permalink, data["title"], data["date"], data["date"], data["description"], data["category"]))
			} else {
				template := `<li>
					<article>
						<h3 class="post-title"><a href="%s">%s</a></h3>
						<p class="post-date"><time datetime="%s">%s</time></p>
						<p class="post-description">%s</p>
						<div class="post-category">%s</div>
					</article>
				</li>`
				postList = append(postList, fmt.Sprintf(template, permalink, data["title"], data["date"], data["date"], data["description"], data["category"]))
			}
		}
		postList = append(postList, "</ul>")
	}
	htmlPostList := strings.Join(postList, "")

	// public/index.html 처리
	type IndexPageTemplateData struct {
		PostList template.HTML
	}
	data := IndexPageTemplateData{
		PostList: template.HTML(htmlPostList),
	}

	tmpl, err := template.ParseFiles("./layout/index.html")
	if err != nil {
		fmt.Printf("layout/index.html 템플릿 파일 파싱 실패\n")
	}

	outputFile, err := os.Create("public/index.html") // 파일이 없으면 새로 생성. 파일이 이미 있으면 초기화.
	if err != nil {
		fmt.Printf("출력 파일 생성 실패\n")
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, data); err != nil {
		fmt.Printf("템플릿 실행 실패\n")
	}
	fmt.Printf("성공: public/index.html 파일 생성.\n")

	// *** post 처리
	// public/post 디렉토리 생성
	var postDirPath string = "public/post"
	if err := os.MkdirAll(postDirPath, 0755); err != nil {
		fmt.Printf("public/post 디렉토리 생성 실패\n")
	} else {
		fmt.Printf("성공: public/post 디렉토리 생성\n")
	}

	md := goldmark.New() // goldmark

	layoutFile := "layout/post.html"
	tmplPost, err := template.ParseFiles(layoutFile)
	if err != nil {
		fmt.Printf("layout/post.html 템플릿 파싱 실패\n")
	}

	type PostPageTemplateData struct {
		Title       string
		Date        string
		Category    []string
		Description string
		URL         string
		Content     template.HTML
	}

	for _, data := range postsData {
		title, _ := data["title"].(string)
		date, _ := data["date"].(string)
		category, _ := data["category"].([]string)
		description, _ := data["description"].(string)
		outputFilePath, _ := data["outputFilePath"].(string)
		permalink := filepath.Join("public", "post", outputFilePath)
		url := filepath.Join("post", outputFilePath)
		sourceFilePath, _ := data["sourceFilePath"].(string)

		// goldmark
		sourceBytes, err := os.ReadFile(sourceFilePath)
		if err != nil {
			fmt.Printf("파일 읽기 실패\n")
			continue
		}
		var bodyBytes []byte
		if bytes.HasPrefix(sourceBytes, []byte("---")) {
			parts := bytes.SplitN(sourceBytes, []byte("---"), 3)
			if len(parts) >= 3 {
				bodyBytes = parts[2]
			} else {
				bodyBytes = sourceBytes
			}
		} else {
			bodyBytes = sourceBytes
		}
		var contentBuf bytes.Buffer
		context := parser.NewContext()
		if err := md.Convert(bodyBytes, &contentBuf, parser.WithContext(context)); err != nil {
			fmt.Printf("Markdown 변환 실패\n")
			continue
		}

		outputFile, err := os.Create(permalink)
		if err != nil {
			fmt.Printf("출력 파일 생성 실패\n")
			continue
		}
		defer outputFile.Close()

		page := PostPageTemplateData{
			Title:       title,
			Date:        date,
			Category:    category,
			Description: description,
			URL:         url,
			Content:     template.HTML(contentBuf.String()),
		}

		if err := tmplPost.Execute(outputFile, page); err != nil {
			fmt.Printf("템플릿 실행 실패\n")
			continue
		}

		fmt.Printf("성공: %s 파일 생성\n", outputFilePath)
	}

	// TODO: assets은 content에서 관리한다.

	// TODO: public/index.html Category 처리하기
}
