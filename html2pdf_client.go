package dgpdf

import (
	"bytes"
	"errors"
	pdf "github.com/adrg/go-wkhtmltopdf"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"html/template"
	"os"
	"sync"
	"time"
)

var (
	initOnce   sync.Once
	semaphore  = make(chan struct{}, 1)
	TimeoutErr = errors.New("get pdfconv engine timeout")
)

func ExportPdfFileByHtmlTemplateFile(ctx *dgctx.DgContext, htmlTemplateFile string, data any, filePath string) error {
	tmpl, err := template.ParseFiles(htmlTemplateFile)
	if err != nil {
		dglogger.Errorf(ctx, "parse html template error, %s, %v", htmlTemplateFile, err)
		return err
	}

	return ExportPdfFileByHtmlTemplate(ctx, tmpl, data, filePath)
}

func ExportPdfFileByHtmlTemplateString(ctx *dgctx.DgContext, htmlTemplateString string, data any, filePath string) error {
	tmpl, err := template.New("html").Parse(htmlTemplateString)
	if err != nil {
		dglogger.Errorf(ctx, "parse html template error, \n%s, \n%v", htmlTemplateString, err)
		return err
	}

	return ExportPdfFileByHtmlTemplate(ctx, tmpl, data, filePath)
}

func ExportPdfFileByHtmlTemplate(ctx *dgctx.DgContext, tmpl *template.Template, data any, filePath string) error {
	var source *bytes.Buffer
	err := tmpl.Execute(source, data)
	if err != nil {
		dglogger.Errorf(ctx, "execute html template parse error: %v", err)
		return err
	}

	initOnce.Do(func() {
		err = pdf.Init()
	})
	got := false
	select {
	case semaphore <- struct{}{}:
		got = true
	case <-time.After(2 * time.Second):
	}
	if !got {
		return TimeoutErr
	}
	defer func() {
		<-semaphore
	}()

	object, err := pdf.NewObjectFromReader(source)
	if err != nil {
		return err
	}

	converter, err := pdf.NewConverter()
	if err != nil {
		return err
	}
	defer converter.Destroy()
	converter.Add(object)
	converter.PaperSize = pdf.A4
	converter.Orientation = pdf.Landscape
	converter.MarginTop = "1cm"
	converter.MarginBottom = "1cm"
	converter.MarginLeft = "10mm"
	converter.MarginRight = "10mm"

	pdfFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer pdfFile.Close()

	err = converter.Run(pdfFile)
	if err != nil {
		return err
	}

	return nil
}
