package excel

import (
	"testing"
	"fmt"
)

func TestCreateExcel(t *testing.T) {
	file, err := CreateExcel(&ExConfig{
		Sheets: []*SheetConfig{
			{
				SheetName: "sheet111",
				Data: [][]interface{}{
					{"姓名", "性别", "年龄"},
					{
						"xx同学",
						"女",
						18,
					},
					{
						"xx同学",
						"男",
						88,
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	file.Save("test.xlsx")
}

func TestCreateExcelToLocal(t *testing.T) {
	err := CreateExcelToLocal(&ExConfig{
		FileName: "test2.xlsx",
		Sheets: []*SheetConfig{
			{
				SheetName: "sheet111",
				Data: [][]interface{}{
					{"姓名", "性别", "年龄"},
					{
						"xx同学",
						"女",
						18,
					},
					{
						"xx同学",
						"男",
						88,
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}