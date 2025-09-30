package main

import (
	"blog/lib"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// Refactor TODO:
/*
0. CI/CD
	1. pnpm run build 하면 build.sh 먼저 실행 및 go run .
	2. pnpm run build:prod 하면 build.sh 먼저 실행 및 cross-env APP_ENV=production go run .

1. data 처리:
	1. build.sh:
		1. unix 명령어로 public 디렉토리 초기화(없으면 생성, 있으면 초기화)
			DOCS: public 디렉토리 자체는 이제 main.go에서 처리하지 않음(즉, 생성 및 초기화 안함)
		2. sass 명령어로 public/styles 디렉토리 생성 및 scss 컴파일
		3. esbuild 명령어로 public/js 디렉토리 생성 및 js 컴파일
		4. unix 명령어로 .nojekyll 파일 생성
		5. unix 명령어로 robots.txt 파일 복사해서 붙여넣기
		6. content/assets 속 모든 파일 public/assets에 붙여넣기
	2. main.go:
		1. 포스트 관련 처리
			1. PostsData 처리
				- title string
				- date string
				- description string
				- category []string
				- fixedPost bool
				- fixedCategory bool
				- sourceFilePath string <- 실제 블로그가 저장된 위치
					e.g. content/tgpl/go 정리법
				- slug <- 그냥 가공된 PATH로 활용될 이름임
					e.g. go-정리법
					-> 여기에 .html, post/, public/ 붙여서 활용하기
			2. post.html 생성
				href: posts/category-foo, /, about, posts
				DOCS: href는 해당 html에 포함된 hyperlink임
			3. posts.html 생성
				href: posts/category-foo, /, about, posts
			4. posts/category-foo.html 생성
				href: /, about, posts
		2. 블로그 관련 처리
			1. index.html 생성
			2. about.html 생성
2.

*/

func main() {
	// Get env
	var appEnv string = os.Getenv("APP_ENV")
	_ = appEnv

	// Set goldmark
	var md = goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // markdown에서 html tag 사용할 수 있게 활성화함
		),
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					// chromahtml.WithLineNumbers(true),
					chromahtml.WithClasses(true),
				),
			),
		),
	)
	_ = md

	// PostsData
	fmt.Println()
	fmt.Println("-- PostsData 처리 --")
	var contentDirPath string = "./content"
	var contentFilePaths []string = lib.GetFilePaths(contentDirPath)
	_ = contentFilePaths

	var postsData []map[string]any
	var postsDataByCategory = make(map[string][]map[string]any)
	_ = postsData
	_ = postsDataByCategory

	type PostFrontMatter struct {
		Date          string   `yaml:"date"`
		Desc          string   `yaml:"desc"`
		Category      []string `yaml:"category"`
		Published     bool     `yaml:"published"`
		FixedPost     bool     `yaml:"fixedPost"`
		FixedCategory bool     `yaml:"fixedCategory"`
	}

	for _, path := range contentFilePaths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("error: %s - 파일 열기 실패\n", path)
			continue
		}
		defer file.Close()
		buffer := make([]byte, 4096)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("error: %s - 파일 읽기 실패\n", path)
			continue
		}
		content := buffer[:n]
		if !bytes.HasPrefix(content, []byte("---")) {
			fmt.Printf("error: %s - 프론트 매터를 찾을 수 없습니다\n", path)
			continue
		}
		parts := bytes.SplitN(content, []byte("---"), 3)
		if len(parts) < 3 {
			fmt.Printf("error: %s - 프론트 매터의 끝을 찾을 수 없습니다\n", path)
			continue
		}
		frontMatterBytes := parts[1]
		var fm PostFrontMatter
		if err := yaml.Unmarshal(frontMatterBytes, &fm); err != nil {
			fmt.Printf("error: %s - YAML 파싱 실패\n", path)
			continue
		}

		// published된 것만 PostsData에 추가한다
		if fm.Published {
			// title
			var title string = path[strings.LastIndex(path, "/")+1 : strings.LastIndex(path, ".")]

			// URL
			var slug string = lib.SlugifyPath(title)

			// description
			var description string
			if fm.Desc != "" {
				description = fm.Desc
			} else {
				bodyContent := bytes.TrimSpace(parts[2])
				plainText, err := lib.MarkdownToPlainText(bodyContent)
				if err != nil {
					fmt.Printf("error: %s - Markdown 변환 실패\n", path)
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

			// PostsData 처리
			var data = map[string]any{
				"title":          title,
				"date":           fm.Date,
				"description":    description,
				"category":       fm.Category,
				"fixedPost":      fm.FixedCategory,
				"fixedCategory":  fm.FixedCategory,
				"sourceFilePath": path,
				"slug":           slug,
			}

			postsData = append(postsData, data)

			for _, category := range fm.Category {
				postsDataByCategory[category] = append(postsDataByCategory[category], data)
			}
		}
	}

	// Post
	fmt.Println()
	fmt.Println("-- Post 처리 --")
	var postDirPath string = "public/post"
	if err := os.MkdirAll(postDirPath, 0755); err != nil {
		fmt.Printf("error: public/post - 디렉토리 생성 실패\n")
	}

	layoutFile := "layout/post.html"
	postTemplate, err := template.ParseFiles(layoutFile)
	if err != nil {
		fmt.Printf("error: layout/post.html - 템플릿 파싱 실패\n")
	}

	type PostTemplateData struct {
		IsProduction  bool
		Title         string
		Date          string
		FormattedDate string
		Category      []string
		Description   string
		Permalink     string
		Content       template.HTML
		CurrentURL    string // for nav tag
	}

	for _, data := range postsData {
		// Post 데이터 처리
		title, _ := data["title"].(string)
		date, _ := data["date"].(string)
		category, _ := data["category"].([]string)
		description, _ := data["description"].(string)
		sourceFilePath, _ := data["sourceFilePath"].(string) // content에 저장된 파일명
		slug, _ := data["slug"].(string)
		publicPath := filepath.Join("public", "post", fmt.Sprintf("%s.html", slug)) // public에 저장된 파일명
		var permalink string                                                        // 블로그의 고유 링크임
		if appEnv == "production" {
			permalink = filepath.ToSlash(filepath.Join("post", slug))
		} else {
			permalink = filepath.ToSlash(filepath.Join("post", fmt.Sprintf("%s.html", slug)))
		}
		formattedDate, err := lib.FormatDateKorean(date) // yyyy년 mm월 dd일
		if err != nil {
			fmt.Printf("error: %s - 날짜 변환 실패\n", sourceFilePath)
			return
		}

		// Post 생성
		sourceBytes, err := os.ReadFile(sourceFilePath)
		if err != nil {
			fmt.Printf("error: %s - 파일 읽기 실패\n", sourceFilePath)
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
			fmt.Printf("error: %s - Markdown 변환 실패\n", sourceFilePath)
			continue
		}

		outputFile, err := os.Create(publicPath)
		if err != nil {
			fmt.Printf("error: %s - 출력 파일 생성 실패\n", sourceFilePath)
			continue
		}
		defer outputFile.Close()

		page := PostTemplateData{
			IsProduction:  appEnv == "production",
			Title:         title,
			Date:          date,
			FormattedDate: formattedDate,
			Category:      category,
			Description:   description,
			Permalink:     permalink,
			Content:       template.HTML(contentBuf.String()),
			CurrentURL:    "/posts",
		}

		if err := postTemplate.Execute(outputFile, page); err != nil {
			fmt.Printf("error: layoutpost.html - 템플릿 실행 실패\n")
			continue
		}

		fmt.Printf("sucess: %s 파일 생성\n", publicPath)
	}

	// Post list 처리
	fmt.Println()
	fmt.Println("-- Post List 처리 --")

}

func isKorean(r rune) bool { return unicode.Is(unicode.Hangul, r) }

func compareKoreanEnglish(s1, s2 string) bool {
	isKorean1 := isKorean([]rune(s1)[0])
	isKorean2 := isKorean([]rune(s2)[0])
	if isKorean1 != isKorean2 {
		return isKorean1 // 한글이 true이므로 앞으로 온다
	}
	return s1 < s2
}
