package ocr

import (
	"github.com/otiai10/gosseract"
)

type IOCR interface {
	ExtractText()
	Close()
}

type OCRInstance struct {
	c *gosseract.Client
}

func New() OCRInstance {
	ocr := &OCRInstance{}
	ocr.c = gosseract.NewClient()
	return ocr
}

func (ocr *OCRInstance) ExtractText() {

}

func (ocr *OCRInstance) Close() {
	defer ocr.c.Close()
}
