package pdf

import "github.com/SebastiaanKlippert/go-wkhtmltopdf"

func ExportPdf(url string, filePath string) (error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}
	pdfg.AddPage(wkhtmltopdf.NewPage(url))
	err = pdfg.Create()
	if err != nil {
		return err
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(filePath)
	if err != nil {
		return err
	}
	return nil
}
