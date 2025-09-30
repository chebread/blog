#!/bin/bash

set -e

echo "-- 빌드 시작 --"

echo "1. public 디렉토리 초기화"
mkdir -p public
find public/ -mindepth 1 -delete

echo "2. sass 컴파일"
pnpm sass layout/styles:public/styles --style=compressed

echo "3. esbuild 컴파일"
pnpm esbuild layout/js/main.ts --bundle --outfile=public/js/main.js --minify --target=es2017

echo "4. .nojekyll 파일 생성"
touch public/.nojekyll

echo "5. robots.txt 파일 복사"
cp layout/robots.txt public/robots.txt

echo "6. assets 디렉토리 복사"
cp -r assets/ public/assets

if [ "$1" = "production" ]; then
  echo "7. 프로덕션 모드로 go 실행"
  pnpm cross-env APP_ENV=production go run .
else
  echo "7. 개발 모드로 go 실행"
  go run .
fi

echo "-- 빌드 완료 --"