package main

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	dgocr "github.com/darwinOrg/go-ppocr"
	"github.com/signintech/gopdf"
	"gopkg.in/gographics/imagick.v3/imagick"
	"strconv"
	"strings"
	"testing"
)

func TestDrawRect(t *testing.T) {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	ctx := &dgctx.DgContext{TraceId: "123"}
	mw.ReadImage("1.pdf")
	var pages = int(mw.GetNumberImages())

	cw := imagick.NewPixelWand()
	defer cw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()

	cw.SetColor("red")
	dw.SetStrokeColor(cw)

	cw.SetAlpha(0)
	dw.SetFillColor(cw)

	dw.SetStrokeWidth(1)
	dw.SetStrokeAntialias(true)

	// 创建PDF文档对象
	p := gopdf.GoPdf{}
	started := false

	for i := 0; i < pages; i++ {
		// This being the page offset
		mw.SetIteratorIndex(i)
		mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE)

		mw.SetImageFormat("jpg")
		mw.SetImageCompression(imagick.COMPRESSION_JPEG)

		pWidth := mw.GetImageWidth()
		pHeight := mw.GetImageHeight()
		mw.ThumbnailImage(pWidth, pHeight)

		if !started {
			// 设置页面大小和边距
			p.Start(gopdf.Config{
				PageSize: gopdf.Rect{
					W: float64(pWidth),
					H: float64(pHeight),
				},
				Protection: gopdf.PDFProtectionConfig{},
			})
			started = true
		}

		imageFile := strconv.Itoa(i+1) + ".jpg"
		mw.WriteImage(imageFile)

		keywords := []string{"Java", "MySQL", "15888888888", "王者荣耀"}

		textRects, _ := dgocr.OcrImageFile(ctx, imageFile)
		var rects []*dgocr.Rect

		p.AddPage()
		p.SetStrokeColor(255, 0, 0)
		p.SetLineWidth(2)
		p.SetFillColor(0, 255, 0)

		for _, tr := range textRects {
			words := tr.Text
			lowerWords := strings.ToLower(words)
			var wordsWidth float64
			leftTopX, leftTopY, rightBottomX, rightBottomY := tr.Rect.LeftTopX, tr.Rect.LeftTopY, tr.Rect.RightBottomX, tr.Rect.RightBottomY

			for _, keyword := range keywords {
				lowerKeyword := strings.ToLower(keyword)
				if strings.Contains(lowerWords, lowerKeyword) {
					keywordIndex := strings.Index(lowerWords, lowerKeyword)
					preWords := words[0:keywordIndex]

					var preWordsWidth float64
					if preWords != "" {
						preWordsMetric := mw.QueryFontMetrics(dw, preWords)
						preWordsWidth = preWordsMetric.TextWidth
					}

					keywordMetrics := mw.QueryFontMetrics(dw, keyword)
					keywordWidth := keywordMetrics.TextWidth

					if wordsWidth == 0 {
						wordsMetric := mw.QueryFontMetrics(dw, words)
						wordsWidth = wordsMetric.TextWidth
					}

					dw.Rectangle(leftTopX+(rightBottomX-leftTopX)*(preWordsWidth/wordsWidth), leftTopY,
						leftTopX+(rightBottomX-leftTopX)*((preWordsWidth+keywordWidth)/wordsWidth), rightBottomY)
					rects = append(rects, &dgocr.Rect{
						LeftTopX:     leftTopX + (rightBottomX-leftTopX)*(preWordsWidth/wordsWidth),
						LeftTopY:     leftTopY,
						RightBottomX: leftTopX + (rightBottomX-leftTopX)*((preWordsWidth+keywordWidth)/wordsWidth),
						RightBottomY: rightBottomY,
					})
				}
			}
		}

		mw.SetImageFormat("jpg")
		mw.SetImageCompression(imagick.COMPRESSION_JPEG)

		if err := mw.DrawImage(dw); err != nil {
			dglogger.Errorf(ctx, "DrawImage error：%v", err)
		}

		if err := mw.WriteImage(imageFile); err != nil {
			dglogger.Errorf(ctx, "WriteImage error：%v", err)
		}

		for _, rect := range rects {
			p.Rectangle(rect.LeftTopX, rect.LeftTopY, rect.RightBottomX, rect.RightBottomY, "DF", 3, 10)
		}
	}

	// 保存PDF文件到本地磁盘
	p.WritePdf("example.pdf")
}
