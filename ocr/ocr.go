package ocr

import (
	"github.com/otiai10/gosseract"
)

//IOCR ocr interface
type IOCR interface {
	ExtractText()
}

type ocrInstance struct {
}

//New instance
func New() IOCR {
	ocr := &ocrInstance{}
	return ocr
}

//ExtractText extract text from image
func (ocr *ocrInstance) ExtractText(imgPath string) (string, error) {
	c = gosseract.NewClient()
	defer c.Close()

	var (
		err  error
		text string
	)

	err = c.SetImage(imgPath)
	if err != nil {
		return "", err
	}

	text, err = c.Text()
	if err != nil {
		return text, err
	}

	return text, err
}

