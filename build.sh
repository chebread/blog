#!/bin/bash

set -e

echo "-- 빌드 시작 --"

echo "public 디렉토리 초기화"
mkdir -p public
find public/ -mindepth 1 -delete

echo "sass 컴파일"
pnpm sass layout/styles:public/styles --style=compressed

echo "esbuild 컴파일"
pnpm esbuild layout/js/main.ts --bundle --outfile=public/js/main.js --minify --target=es2017

echo ".nojekyll 파일 생성"
touch public/.nojekyll

echo "robots.txt 파일 복사"
cp layout/robots.txt public/robots.txt

echo "assets 디렉토리 복사"
cp -r assets/ public/assets

echo "favicons 디렉토리 복사"
cp -r layout/favicons public/favicons

if [ "$1" = "production" ]; then
  echo "프로덕션 모드로 go 실행"
  pnpm cross-env APP_ENV=production go run .
else
  echo "개발 모드로 go 실행"
  go run .
fi

echo "-- 빌드 완료 --"