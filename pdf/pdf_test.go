package pdf

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHtml2Pdf(t *testing.T){
	url := "https://segmentfault.com/a/1190000015749615?utm_source=tag-newest"
	filePath := "/home/qvdev/桌面test.pdf"
	err := ExportPdf(url, filePath)
	assert.NoError(t, err)
}