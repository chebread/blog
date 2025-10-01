package main

import (
	"blog/lib"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
				- fixed bool
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
		Date      string   `yaml:"date"`
		Desc      string   `yaml:"desc"`
		Category  []string `yaml:"category"`
		Published bool     `yaml:"published"`
		Fixed     bool     `yaml:"fixed"`
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
				"fixed":          fm.Fixed,
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

	type PostTmplData struct {
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

		postTmplData := PostTmplData{
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

		if err := postTemplate.Execute(outputFile, postTmplData); err != nil {
			fmt.Printf("error: layoutpost.html - 템플릿 실행 실패\n")
			continue
		}

		fmt.Printf("sucess: %s 파일 생성\n", publicPath)
	}

	// index.html 처리
	fmt.Println()
	fmt.Println("-- index.html 처리 --")

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
		IsProduction bool
		CurrentURL   string
		Content      template.HTML
	}
	homePageData := HomePageData{
		IsProduction: appEnv == "production",
		CurrentURL:   "/",
		Content:      template.HTML(homeContentBuf.String()),
	}

	if err := tmplHome.Execute(outputHomeFile, homePageData); err != nil {
		fmt.Printf("템플릿 실행 실패: %v\n", err)
		return
	}

	fmt.Printf("성공: %s 파일 생성\n", destHomeFile)

	// about.html 처리
	fmt.Println()
	fmt.Println("-- abou.html 처리 --")

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
		IsProduction bool
		CurrentURL   string
		Content      template.HTML
	}
	aboutPageData := AboutPageData{
		IsProduction: appEnv == "production",
		CurrentURL:   "/about",
		Content:      template.HTML(aboutContentBuf.String()),
	}

	if err := tmplabout.Execute(outputAboutFile, aboutPageData); err != nil {
		fmt.Printf("템플릿 실행 실패: %v\n", err)
		return
	}

	fmt.Printf("성공: %s 파일 생성\n", destAboutFile)

	// Post list 처리
	fmt.Println()
	fmt.Println("-- Post List 처리 --")

	var postList []string // postListHtml strings.Builder 방식으로 바꾸기.

	var categories []string // map은 순서가 보장되지 않으니까 이런 방식으로 순서화한다.
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
		var categoryLink string
		if appEnv == "production" {
			categoryLink = filepath.ToSlash(filepath.Join("posts", lib.SlugifyPath(category)))
		} else {
			categoryLink = filepath.ToSlash(filepath.Join("posts", fmt.Sprintf("%s.html", lib.SlugifyPath(category))))
		}
		_ = categoryLink

		posts := postsDataByCategory[category]
		postList = append(postList, "<section class=\"category-group\">")
		postList = append(postList, fmt.Sprintf("<h2 class=\"category-group-title\"><a href=\"%s\">[%s]</a></h2>", categoryLink, category))
		postList = append(postList, "<ul class=\"category-group-list\">")

		// 정렬: fixed, date
		sort.Slice(posts, func(i, j int) bool {
			postI := posts[i]
			postJ := posts[j]

			fixedI, _ := postI["fixed"].(bool)
			fixedJ, _ := postJ["fixed"].(bool)

			if fixedI != fixedJ {
				return fixedI
			}

			if fixedI {
				titleI, _ := postI["title"].(string)
				titleJ, _ := postJ["title"].(string)

				isKoreanI := lib.IsKorean(titleI)
				isKoreanJ := lib.IsKorean(titleJ)

				if isKoreanI != isKoreanJ {
					return isKoreanI
				}

				return titleI < titleJ
			}

			dateI, _ := postI["date"].(string)
			dateJ, _ := postJ["date"].(string)
			return dateI > dateJ
		})

		const maxPostsToShow = 3
		postsToDisplay := posts
		needsMoreLink := false

		if len(posts) > maxPostsToShow {
			postsToDisplay = posts[:maxPostsToShow]
			needsMoreLink = true
		}

		for _, data := range postsToDisplay {
			isFixed, _ := data["fixed"].(bool)
			title, _ := data["title"].(string)
			date, _ := data["date"].(string) // yyyy-mm-dd
			description, _ := data["description"].(string)
			// categorySlice, _ := data["category"].([]string)
			slug, _ := data["slug"].(string)

			var permalink string // 블로그 링크는 Production에서는 .html가 없음
			if appEnv == "production" {
				permalink = filepath.ToSlash(filepath.Join("post", slug))
			} else {
				permalink = filepath.ToSlash(filepath.Join("post", fmt.Sprintf("%s.html", slug)))
			}

			formattedDate, err := lib.FormatDateKorean(date) // yyyy년 mm월 dd일
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
                    </article>
                </li>`
				postList = append(postList, fmt.Sprintf(fixedTemplate, permalink, title, date, formattedDate, description))
			} else {
				template := `<li>
                    <article class="post-item">
                        <h3 class="post-item-title"><a href="%s">%s</a></h3>
                        <p class="post-item-date"><time datetime="%s">%s</time></p>
                        <p class="post-item-description">%s</p>
                    </article>
                </li>`
				postList = append(postList, fmt.Sprintf(template, permalink, title, date, formattedDate, description))
			}
		}

		if needsMoreLink {
			var moreLinkURL string // 블로그 링크는 Production에서는 .html가 없음
			if appEnv == "production" {
				moreLinkURL = filepath.ToSlash(filepath.Join("posts", lib.SlugifyPath((category))))
			} else {
				moreLinkURL = filepath.ToSlash(filepath.Join("posts", fmt.Sprintf("%s.html", lib.SlugifyPath(category))))
			}
			moreLinkHTML := fmt.Sprintf(`<li><article class="post-more-link"><a href="%s">[%s] 더보기</a></article></li>`, moreLinkURL, category)
			postList = append(postList, moreLinkHTML)
		}

		postList = append(postList, "</ul>")
		postList = append(postList, "</section>")
	}
	htmlPostList := strings.Join(postList, "")

	type PostsPageTemplateData struct {
		IsProduction bool
		PostList     template.HTML
		CurrentURL   string
	}
	postsPageTemplateData := PostsPageTemplateData{
		IsProduction: appEnv == "production",
		PostList:     template.HTML(htmlPostList),
		CurrentURL:   "/posts",
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

	// Category 처리
	fmt.Println()
	fmt.Println("-- Category별 Post page 생성 --")

	categoryPageTmpl, err := template.ParseFiles("./layout/category.html")
	if err != nil {
		fmt.Printf("layout/category_page.html 템플릿 파일 파싱 실패: %v\n", err)
		return
	}

	categoryPostDir := "public/posts"
	if err := os.MkdirAll(categoryPostDir, 0755); err != nil {
		fmt.Printf("디렉토리 생성 실패: %v\n", err)
		return
	}

	for category, posts := range postsDataByCategory {
		sort.Slice(posts, func(i, j int) bool {
			postI := posts[i]
			postJ := posts[j]
			fixedI, _ := postI["fixed"].(bool)
			fixedJ, _ := postJ["fixed"].(bool)
			if fixedI != fixedJ {
				return fixedI
			}
			if fixedI {
				titleI, _ := postI["title"].(string)
				titleJ, _ := postJ["title"].(string)
				isKoreanI := lib.IsKorean(titleI)
				isKoreanJ := lib.IsKorean(titleJ)
				if isKoreanI != isKoreanJ {
					return isKoreanI
				}
				return titleI < titleJ
			}
			dateI, _ := postI["date"].(string)
			dateJ, _ := postJ["date"].(string)
			return dateI > dateJ
		})

		var postListHtml strings.Builder
		postListHtml.WriteString("<section class=\"category-group\">")
		postListHtml.WriteString(fmt.Sprintf("<h2 class=\"category-group-title\">[%s]</h2>", category))
		postListHtml.WriteString("<ul class=\"category-group-list\">")

		for _, data := range posts {
			isFixed, _ := data["fixed"].(bool)
			title, _ := data["title"].(string)
			date, _ := data["date"].(string)
			description, _ := data["description"].(string)
			slug, _ := data["slug"].(string)

			var permalink string
			if appEnv == "production" {
				permalink = filepath.ToSlash(filepath.Join("/", "post", slug))
			} else {
				permalink = filepath.ToSlash(filepath.Join("/", "post", fmt.Sprintf("%s.html", slug)))
			}

			formattedDate, _ := lib.FormatDateKorean(date)

			var templateString string
			if isFixed {
				templateString = `<li><article class="post-item"><h3 class="post-item-title"><a href="%s">[고정됨] %s</a></h3><p class="post-item-date"><time datetime="%s">%s</time></p><p class="post-item-description">%s</p></article></li>`
				postListHtml.WriteString(fmt.Sprintf(templateString, permalink, title, date, formattedDate, description))
			} else {
				templateString = `<li><article class="post-item"><h3 class="post-item-title"><a href="%s">%s</a></h3><p class="post-item-date"><time datetime="%s">%s</time></p><p class="post-item-description">%s</p></article></li>`
				postListHtml.WriteString(fmt.Sprintf(templateString, permalink, title, date, formattedDate, description))
			}
		}

		var backLinkURL string
		if appEnv == "production" {
			backLinkURL = "/posts"
		} else {
			backLinkURL = "/posts.html"
		}
		backButtonHTML := fmt.Sprintf(`<li><article class="back-link"><a href="%s">돌아가기</a></article></li>`, backLinkURL)
		postListHtml.WriteString(backButtonHTML)

		postListHtml.WriteString("</ul>")
		postListHtml.WriteString("</section>")

		type CategoryPageTemplateData struct {
			IsProduction bool
			CategoryName string
			PostList     template.HTML
			CurrentURL   string
		}

		templateData := CategoryPageTemplateData{
			IsProduction: appEnv == "production",
			CategoryName: category,
			PostList:     template.HTML(postListHtml.String()),
			CurrentURL:   "/posts",
		}

		var outputFilePath string
		if appEnv == "production" {
			outputFilePath = filepath.Join(categoryPostDir, lib.SlugifyPath(category), "index.html")
			os.MkdirAll(filepath.Join(categoryPostDir, lib.SlugifyPath(category)), 0755)
		} else {
			outputFilePath = filepath.Join(categoryPostDir, fmt.Sprintf("%s.html", lib.SlugifyPath(category)))
		}

		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			fmt.Printf("출력 파일 생성 실패 (%s): %v\n", outputFilePath, err)
			continue
		}
		defer outputFile.Close()

		if err := categoryPageTmpl.Execute(outputFile, templateData); err != nil {
			fmt.Printf("템플릿 실행 실패 (%s): %v\n", outputFilePath, err)
		}

		fmt.Printf("성공: %s 파일 생성\n", outputFilePath)
	}
}
