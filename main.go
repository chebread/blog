package main

import (
	"blog/lib"
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// TODO: post에 Thumnail 추가 -> header hero 처럼

// TODO: postlist, category list의 a tag에서만 퍼센트 인코딩해서 전달하면 된다.

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
		goldmark.WithParserOptions(
			parser.WithAttribute(), // 수동 id 지정 기능
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

	type CategoryInfo struct {
		Name string
		URL  string
	}

	type PostTmplData struct {
		IsProduction  bool
		Title         string
		Date          string
		FormattedDate string
		// Category      []string
		Description string
		Permalink   string
		Content     template.HTML
		CurrentURL  string // for nav tag
		Categories  []CategoryInfo
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

		var categoriesData []CategoryInfo
		for _, c := range category {
			var linkURL string
			slugifiedPath := lib.SlugifyPath(c)

			if appEnv == "production" {
				linkURL = filepath.ToSlash(filepath.Join("/", "posts", slugifiedPath))
			} else {
				linkURL = filepath.ToSlash(filepath.Join("/", "posts", fmt.Sprintf("%s.html", slugifiedPath)))
			}

			categoriesData = append(categoriesData, CategoryInfo{
				Name: c,       // "Go"
				URL:  linkURL, // "/posts/go.html"
			})
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
			// Category:      category,
			Description: description,
			Permalink:   permalink,
			Content:     template.HTML(contentBuf.String()),
			CurrentURL:  "/posts",
			Categories:  categoriesData,
		}

		if err := postTemplate.Execute(outputFile, postTmplData); err != nil {
			fmt.Printf("error: layoutpost.html - 템플릿 실행 실패\n")
			continue
		}

		fmt.Printf("sucess: %s 파일 생성\n", publicPath)
	}

	// Post list 처리
	fmt.Println()
	fmt.Println("-- Post List 처리 --")

	var postList []string // postListHtml strings.Builder 방식으로 바꾸기.

	var categories []string // map은 순서가 보장되지 않으니까 이런 방식으로 순서화한다.
	for category := range postsDataByCategory {
		categories = append(categories, category)
	}

	// 목록 Sort
	sort.Slice(categories, func(i, j int) bool {
		return lib.CompareStrings(categories[i], categories[j])
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

		// fixed 정렬
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
				return lib.CompareStrings(titleI, titleJ)
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
			encodedSlug := url.PathEscape(slug)

			var permalink string // 블로그 링크는 Production에서는 .html가 없음
			if appEnv == "production" {
				permalink = filepath.ToSlash(filepath.Join("post", encodedSlug))
			} else {
				permalink = filepath.ToSlash(filepath.Join("post", fmt.Sprintf("%s.html", encodedSlug)))
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
			moreLinkHTML := fmt.Sprintf(`<li><article class="post-more-link"><a href="%s">더보기...</a></article></li>`, moreLinkURL)
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

	// 정렬
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
				return lib.CompareStrings(titleI, titleJ)
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
			encodedSlug := url.PathEscape(slug)

			var permalink string
			if appEnv == "production" {
				permalink = filepath.ToSlash(filepath.Join("/", "post", encodedSlug))
			} else {
				permalink = filepath.ToSlash(filepath.Join("/", "post", fmt.Sprintf("%s.html", encodedSlug)))
			}

			formattedDate, _ := lib.FormatDateKorean(date)

			var templateString string
			if isFixed {
				templateString = `<li>
					<article class="category-item">
						<h3 class="category-item-title"><a href="%s">[고정됨] %s</a></h3>
						<p class="category-item-date"><time datetime="%s">%s</time></p>
						<p class="category-item-description">%s</p>
					</article>
				</li>`
				postListHtml.WriteString(fmt.Sprintf(templateString, permalink, title, date, formattedDate, description))
			} else {
				templateString = `<li>
					<article class="category-item">
						<h3 class="category-item-title"><a href="%s">%s</a></h3>
						<p class="category-item-date"><time datetime="%s">%s</time></p>
						<p class="category-item-description">%s</p>
					</article>
				</li>`
				postListHtml.WriteString(fmt.Sprintf(templateString, permalink, title, date, formattedDate, description))
			}
		}

		backButtonHTML :=
			`<li>
				<article class="back-link">
					<a href="#" id="back">돌아가기...</a>
				</article>
			</li>`
			// onclick="history.back(); return false;"
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

		var outputFilePath = filepath.Join(categoryPostDir, fmt.Sprintf("%s.html", lib.SlugifyPath(category)))

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

	// sitemap.xml 처리
	type SitemapURL struct {
		XMLName    xml.Name `xml:"url"`
		Loc        string   `xml:"loc"`
		LastMod    string   `xml:"lastmod"`
		ChangeFreq string   `xml:"changefreq"`
	}

	type URLSet struct {
		XMLName     xml.Name     `xml:"urlset"`
		Xmlns       string       `xml:"xmlns,attr"`
		XmlnsXSI    string       `xml:"xmlns:xsi,attr"`
		XSILocation string       `xml:"xsi:schemaLocation,attr"`
		URLs        []SitemapURL `xml:"url"`
	}

	fmt.Println()
	fmt.Println("-- Sitemap 생성 --")

	const baseURL = "https://chebread.github.io"

	var urlset = &URLSet{
		Xmlns:       "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsXSI:    "http://www.w3.1/XMLSchema-instance",
		XSILocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
	}

	var today = time.Now().Format("2006-01-02")

	var staticPages = []string{"", "about.html"} // index.html은 경로가 "" 임.
	for _, page := range staticPages {
		urlEntry := SitemapURL{
			Loc:        fmt.Sprintf("%s/%s", baseURL, page),
			LastMod:    today,
			ChangeFreq: "monthly",
		}
		urlset.URLs = append(urlset.URLs, urlEntry)
	}

	for _, data := range postsData {
		var slug, okSlug = data["slug"].(string)
		var date, okDate = data["date"].(string)

		if !okSlug || !okDate {
			continue
		}

		var postURL string
		if appEnv == "production" {
			postURL = fmt.Sprintf("%s/post/%s", baseURL, slug)
		} else {
			postURL = fmt.Sprintf("%s/post/%s.html", baseURL, slug)
		}

		var urlEntry = SitemapURL{
			Loc:        postURL,
			LastMod:    date,     // 포스트의 실제 최종 수정일
			ChangeFreq: "weekly", // 포스트는 비교적 자주 변경될 수 있으므로 weekly로 설정함
		}
		urlset.URLs = append(urlset.URLs, urlEntry)
	}

	xmlBytes, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		fmt.Printf("error: Sitemap XML 변환 실패\n")
		return
	}

	var sitemapContent = []byte(xml.Header + string(xmlBytes))

	var sitemapPath = "public/sitemap.xml"
	err = os.WriteFile(sitemapPath, sitemapContent, 0644)
	if err != nil {
		fmt.Printf("error: %s 파일 생성 실패\n", sitemapPath)
		return
	}

	fmt.Printf("성공: %s 파일 생성\n", sitemapPath)
}
