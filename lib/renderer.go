package lib

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// MarkdownToPlainText는 Markdown 바이트 슬라이스를 받아 순수 텍스트 문자열로 변환합니다.
// 이 함수만 외부로 노출(Export)됩니다.
func MarkdownToPlainText(markdownContent []byte) (string, error) {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithRenderer(renderer.NewRenderer(renderer.WithNodeRenderers(
			util.Prioritized(&plainTextRenderer{}, 1000),
		))),
	)

	if err := md.Convert(markdownContent, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// plainTextRenderer는 mdutil 패키지 내부에서만 사용됩니다. (소문자로 시작)
type plainTextRenderer struct{}

func (r *plainTextRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindCodeSpan, r.renderText)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindHeading, r.renderParagraph)
	// ... (필요에 따라 다른 노드 타입 추가)
}

func (r *plainTextRenderer) renderText(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.Write(n.(*ast.Text).Segment.Value(source))
	}
	return ast.WalkContinue, nil
}

func (r *plainTextRenderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		w.WriteString(" ")
	}
	return ast.WalkContinue, nil
}
