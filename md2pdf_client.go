package dgpdf

import (
	"bytes"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	pdf "github.com/stephenafamo/goldmark-pdf"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"os"
	"text/template"
)

func ExportPdfFileByMarkdownTemplateFile(ctx *dgctx.DgContext, mdTemplateFile string, data any, filePath string) error {
	tmpl, err := template.ParseFiles(mdTemplateFile)
	if err != nil {
		dglogger.Errorf(ctx, "parse markdown template error, %s, %v", mdTemplateFile, err)
		return err
	}

	return ExportPdfFileByMarkdownTemplate(ctx, tmpl, data, filePath)
}

func ExportPdfFileByMarkdownTemplateString(ctx *dgctx.DgContext, mdTemplateString string, data any, filePath string) error {
	tmpl, err := template.New("md").Parse(mdTemplateString)
	if err != nil {
		dglogger.Errorf(ctx, "parse markdown template error, \n%s, \n%v", mdTemplateString, err)
		return err
	}

	return ExportPdfFileByMarkdownTemplate(ctx, tmpl, data, filePath)
}

func ExportPdfFileByMarkdownTemplate(ctx *dgctx.DgContext, tmpl *template.Template, data any, filePath string) error {
	var source bytes.Buffer
	err := tmpl.Execute(&source, data)
	if err != nil {
		dglogger.Errorf(ctx, "execute markdown template parse error: %v", err)
		return err
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRenderer(pdf.New()),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var dest bytes.Buffer
	err = md.Convert(source.Bytes(), &dest)
	if err != nil {
		dglogger.Errorf(ctx, "convert pdf error: %v", err)
		return err
	} else {
		err := os.WriteFile(filePath, dest.Bytes(), os.ModePerm)
		if err != nil {
			dglogger.Errorf(ctx, "write pdf file error: %v", err)
			return err
		}
	}

	return nil
}
