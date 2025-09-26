package main

import (
	"blog/build"
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
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// TODO: IO 공부하기

// TODO: SEO image 추가하기

// TODO: Email 추가하기

// FIX: safari에서 css 이상함. 링크 클릭시 줄어드는 오류 해결하기.

// TODO: 다크모드

func main() {
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // markdown에서 html tag 사용할 수 있게 활성화함
		),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					// chromahtml.WithLineNumbers(true),
					chromahtml.WithClasses(true),
				),
			),
		),
	) // goldmark

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

	// *** 정적 파일 준비
	if err := build.SetupStaticAssets(); err != nil {
		log.Fatalf("정적 파일 준비 중 에러 발생: %v", err)
	}
	fmt.Printf("성공: 정적 파일 준비 완료\n")

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
		postList = append(postList, "<section class=\"category-group\">")
		postList = append(postList, fmt.Sprintf("<h2 class=\"category-group-title\">[%s]</h2>", category))
		postList = append(postList, "<ul class=\"category-group-list\">")

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
			title, _ := data["title"].(string)
			date, _ := data["date"].(string) // yyyy-mm-dd
			description, _ := data["description"].(string)
			category, _ := data["category"].([]string)
			permalink := filepath.Join("post", outputFilePath) // 블로그 링크
			formattedDate, err := lib.FormatDateKorean(date)   // yyyy년 mm월 dd일
			if err != nil {
				fmt.Printf("날짜 변환 실패: %v\n", err)
				return
			}

			if isFixed {
				fixedTemplate := `<li>
					<article class="post-item">
						<h3 class="post-item-title"><a href="%s">[고정됨] %s</a></h3>
						<p class="post-item-date"><time datetime="%s">%s</time></p>
						<p class="post-item-description">%s</p>
						<p class="post-item-category">%s</p>
					</article>
				</li>`
				postList = append(postList, fmt.Sprintf(fixedTemplate, permalink, title, date, formattedDate, description, category))
			} else {
				template := `<li>
					<article class="post-item">
						<h3 class="post-item-title"><a href="%s">%s</a></h3>
						<p class="post-item-date"><time datetime="%s">%s</time></p>
						<p class="post-item-description">%s</p>
						<p class="post-item-category">%s</p>
					</article>
				</li>`
				postList = append(postList, fmt.Sprintf(template, permalink, title, date, formattedDate, description, category))
			}
		}
		postList = append(postList, "</ul>")
		postList = append(postList, "</section>")
	}
	htmlPostList := strings.Join(postList, "")

	// index.html 처리
	homeMdBytes, err := os.ReadFile("content/home/home.md")
	if err != nil {
		fmt.Printf("home.md 파일 읽기 실패: %v\n", err)
		return
	}

	var homeBodyBytes []byte
	if bytes.HasPrefix(homeMdBytes, []byte("---")) {
		parts := bytes.SplitN(homeMdBytes, []byte("---"), 3)
		if len(parts) >= 3 {
			homeBodyBytes = parts[2]
		} else {
			homeBodyBytes = homeMdBytes
		}
	} else {
		homeBodyBytes = homeMdBytes
	}

	var homeContentBuf bytes.Buffer
	if err := md.Convert(homeBodyBytes, &homeContentBuf); err != nil {
		fmt.Printf("home.md 변환 실패: %v\n", err)
		return
	}

	sourceHomeFile := "layout/index.html"
	tmplHome, err := template.ParseFiles(sourceHomeFile)
	if err != nil {
		fmt.Printf("템플릿 파일 파싱 실패: %v\n", err)
		return
	}

	destHomeFile := "public/index.html"
	outputHomeFile, err := os.Create(destHomeFile)
	if err != nil {
		fmt.Printf("출력 파일 생성 실패: %v\n", err)
		return
	}
	defer outputHomeFile.Close()

	type HomePageData struct {
		CurrentURL string
		Content    template.HTML
	}
	homePageData := HomePageData{
		CurrentURL: "/",
		Content:    template.HTML(homeContentBuf.String()),
	}

	if err := tmplHome.Execute(outputHomeFile, homePageData); err != nil {
		fmt.Printf("템플릿 실행 실패: %v\n", err)
		return
	}

	fmt.Printf("성공: %s 파일 생성\n", destHomeFile)

	// about.html 처리
	aboutMdBytes, err := os.ReadFile("content/about/about.md")
	if err != nil {
		fmt.Printf("about.md 파일 읽기 실패: %v\n", err)
		return
	}

	var aboutBodyBytes []byte
	if bytes.HasPrefix(aboutMdBytes, []byte("---")) {
		parts := bytes.SplitN(aboutMdBytes, []byte("---"), 3)
		if len(parts) >= 3 {
			aboutBodyBytes = parts[2]
		} else {
			aboutBodyBytes = aboutMdBytes
		}
	} else {
		aboutBodyBytes = aboutMdBytes
	}

	var aboutContentBuf bytes.Buffer
	if err := md.Convert(aboutBodyBytes, &aboutContentBuf); err != nil {
		fmt.Printf("about.md 변환 실패: %v\n", err)
		return
	}

	sourceAboutFile := "layout/about.html"
	tmplabout, err := template.ParseFiles(sourceAboutFile)
	if err != nil {
		fmt.Printf("템플릿 파일 파싱 실패: %v\n", err)
		return
	}

	destAboutFile := "public/about.html"
	outputAboutFile, err := os.Create(destAboutFile)
	if err != nil {
		fmt.Printf("출력 파일 생성 실패: %v\n", err)
		return
	}
	defer outputAboutFile.Close()

	type AboutPageData struct {
		CurrentURL string
		Content    template.HTML
	}
	aboutPageData := AboutPageData{
		CurrentURL: "/about.html",
		Content:    template.HTML(aboutContentBuf.String()),
	}

	if err := tmplabout.Execute(outputAboutFile, aboutPageData); err != nil {
		fmt.Printf("템플릿 실행 실패: %v\n", err)
		return
	}

	fmt.Printf("성공: %s 파일 생성\n", destAboutFile)

	// *** Posts 처리
	type PostsPageTemplateData struct {
		PostList   template.HTML
		CurrentURL string
	}
	postsPageTemplateData := PostsPageTemplateData{
		PostList:   template.HTML(htmlPostList),
		CurrentURL: "/posts.html",
	}

	tmpl, err := template.ParseFiles("./layout/posts.html")
	if err != nil {
		fmt.Printf("layout/posts.html 템플릿 파일 파싱 실패\n")
	}

	outputFile, err := os.Create("public/posts.html") // 파일이 없으면 새로 생성. 파일이 이미 있으면 초기화.
	if err != nil {
		fmt.Printf("출력 파일 생성 실패\n")
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, postsPageTemplateData); err != nil {
		fmt.Printf("템플릿 실행 실패\n")
	}
	fmt.Printf("성공: public/posts.html 파일 생성\n")

	// *** post 처리
	// public/post 디렉토리 생성
	var postDirPath string = "public/post"
	if err := os.MkdirAll(postDirPath, 0755); err != nil {
		fmt.Printf("public/post 디렉토리 생성 실패\n")
	} else {
		fmt.Printf("성공: public/post 디렉토리 생성\n")
	}

	layoutFile := "layout/post.html"
	tmplPost, err := template.ParseFiles(layoutFile)
	if err != nil {
		fmt.Printf("layout/post.html 템플릿 파싱 실패\n")
	}

	type PostPageTemplateData struct {
		Title         string
		Date          string
		FormattedDate string
		Category      []string
		Description   string
		URL           string
		Content       template.HTML
		CurrentURL    string
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

		formattedDate, err := lib.FormatDateKorean(date) // yyyy년 mm월 dd일
		if err != nil {
			fmt.Printf("날짜 변환 실패: %v\n", err)
			return
		}

		// goldmark
		// TODO: inline code, code block hightlight 설정
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
			Title:         title,
			Date:          date,
			FormattedDate: formattedDate,
			Category:      category,
			Description:   description,
			URL:           url,
			Content:       template.HTML(contentBuf.String()),
			CurrentURL:    "/posts.html",
		}

		if err := tmplPost.Execute(outputFile, page); err != nil {
			fmt.Printf("템플릿 실행 실패\n")
			continue
		}

		fmt.Printf("성공: %s 파일 생성\n", outputFilePath)
	}

}
