package dgpdf

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"os"
	"testing"
)

type inventory struct {
	Name  string
	Count uint
}

var data = []inventory{
	{
		Name:  "name1",
		Count: 1,
	},
	{
		Name:  "name2",
		Count: 2,
	},
}

func TestExportPdfFileByMarkdownTemplateFile(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	err := ExportPdfFileByMarkdownTemplateFile(ctx, "markdown.md", data, "export.pdf")
	if err != nil {
		return
	}
}

func TestExportPdfFileByMarkdownTemplateString(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	bytes, err := os.ReadFile("markdown.md")
	if err != nil {
		return
	}
	err = ExportPdfFileByMarkdownTemplateString(ctx, string(bytes), data, "export.pdf")
	if err != nil {
		return
	}
}
