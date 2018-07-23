package docx

import (
	"github.com/nguyenthenguyen/docx"
)

// 根据模板生成文档，replaces的key是模板中要替换的字符串，value是需要替换成的值
func GenDocxFromTemp(tempFile string, saveToFile string, replaces map[string]string) error {
	r, err := docx.ReadDocxFile(tempFile)
	if err != nil {
		return err
	}
	docx1 := r.Editable()
	for k, v := range replaces {
		// 替换文本正文内容
		docx1.Replace(k, v, -1)
		// 替换页眉内容
		docx1.ReplaceHeader(k, v)
		// 替换页脚内容
		docx1.ReplaceFooter(k, v)
		// 替换链接
		docx1.ReplaceLink(k, v, -1)
	}
	docx1.WriteToFile(saveToFile)
	r.Close()
	return nil
}