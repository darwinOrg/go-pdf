package dgpdf

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"os"
	"testing"
)

func TestExportPdfFileByHtmlTemplateFile(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	err := ExportPdfFileByHtmlTemplateFile(ctx, "template.html", data, "export.pdf")
	if err != nil {
		return
	}
}

func TestExportPdfFileByHtmlTemplateString(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	bytes, err := os.ReadFile("template.html")
	if err != nil {
		return
	}
	err = ExportPdfFileByHtmlTemplateString(ctx, string(bytes), data, "export.pdf")
	if err != nil {
		return
	}
}
