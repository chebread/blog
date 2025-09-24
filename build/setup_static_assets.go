package build

import (
	"blog/lib"
	"fmt"
	"os"
	"path/filepath"
)

func SetupStaticAssets() error {
	// public 디렉토리 삭제 및 생성
	publicDir := "public"
	if err := lib.InitDir(publicDir); err != nil {
		fmt.Printf("public 디렉토리 초기화 실패\n")
	}

	sourceFaviconsDir := "layout/favicons"
	destFaviconsDir := "public/favicons"
	if err := lib.CopyDir(sourceFaviconsDir, destFaviconsDir); err != nil {
		fmt.Printf("layout/favicons 디렉토리 복사 실패\n")
	}
	fmt.Printf("성공: layout/favicons 디렉토리 복사\n")

	// assets은 블로그 글을 위한 이미지... 파일 및 블로그 프로그램을 위한 이미지... 파일을 위한 디렉토리이다
	// TODO: assets 이미지 파일 압축하기
	sourceAssetsDir := "assets"
	destAssetsDir := "public/assets"
	if err := lib.CopyDir(sourceAssetsDir, destAssetsDir); err != nil {
		fmt.Printf("assets 디렉토리 복사 실패\n")
	}
	fmt.Printf("성공: public/assets 디렉토리 복사\n")

	// gp-pages에서 기본적으로 제공하는 404 사용함
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

	source404File := "layout/robots.txt"
	dest404File := "public/robots.txt"
	if err := lib.CopyFile(source404File, dest404File); err != nil {
		fmt.Printf("layout/robots.txt 파일 복사 실패\n")
	}
	fmt.Printf("성공: layout/robots.txt 파일 복사\n")

	return nil
}
