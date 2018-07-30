package pdf

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"fmt"
)

func ExportPdf(url string, filePath string) (error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		fmt.Println("----")
		fmt.Println(err.Error())
		return err
	}
	pdfg.AddPage(wkhtmltopdf.NewPage(url))
	err = pdfg.Create()
	if err != nil {
		fmt.Println("111")
		return err
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(filePath)
	if err != nil {
		fmt.Println("22")
		return err
	}
	return nil
}
