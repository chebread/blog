package lib

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// MarkdownToPlainText는 Markdown 바이트 슬라이스를 받아 순수 텍스트 문자열로 변환합니다.
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

// plainTextRenderer는 패키지 내부에서만 사용됩니다.
type plainTextRenderer struct{}

// RegisterFuncs는 렌더러가 각 Markdown 요소(노드)를 어떻게 처리할지 규칙을 등록합니다.
func (r *plainTextRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// 인라인 요소
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)

	// 블록 요소
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindHeading, r.renderParagraph)
	reg.Register(ast.KindList, r.renderParagraph)
	reg.Register(ast.KindListItem, r.renderParagraph)
	reg.Register(ast.KindBlockquote, r.renderParagraph)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)

	// 무시할 요소
	reg.Register(ast.KindHTMLBlock, r.renderNothing)
	reg.Register(ast.KindRawHTML, r.renderNothing)
}

// --- 각 요소별 처리 함수들 ---

func (r *plainTextRenderer) renderText(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if text, ok := n.(*ast.Text); ok {
			// 텍스트의 실제 내용을 씁니다.
			w.Write(text.Segment.Value(source))

			// 만약 이 텍스트 노드가 소프트 줄 바꿈으로 끝난다면,
			// 띄어쓰기를 보장하기 위해 공백 하나를 추가합니다.
			if text.SoftLineBreak() {
				w.WriteString(" ")
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *plainTextRenderer) renderCodeSpan(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	// CodeSpan은 자식으로 Text 노드를 가지므로, 자식 노드를 계속 순회하도록 합니다.
	// 별도의 타입 단언이 필요 없습니다.
	return ast.WalkContinue, nil
}

func (r *plainTextRenderer) renderLink(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	// 링크는 URL 대신 텍스트만 렌더링합니다. 렌더링이 끝나면 공백을 추가합니다.
	if !entering {
		w.WriteString(" ")
	}
	return ast.WalkContinue, nil
}

func (r *plainTextRenderer) renderImage(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	// 이미지는 URL 대신 alt text만 렌더링합니다.
	if entering {
		if img, ok := n.(*ast.Image); ok {
			w.Write(img.Text(source))
		}
	}
	if !entering {
		w.WriteString(" ")
	}
	return ast.WalkSkipChildren, nil
}

func (r *plainTextRenderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		w.WriteString(" ")
	}
	return ast.WalkContinue, nil
}

func (r *plainTextRenderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			w.Write(line.Value(source))
		}
	}
	if !entering {
		w.WriteString(" ")
	}
	return ast.WalkContinue, nil
}

func (r *plainTextRenderer) renderNothing(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}
