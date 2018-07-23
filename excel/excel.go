package excel

import (
	"github.com/tealeg/xlsx"
	"time"
	"strconv"
	"math/rand"
	"errors"
	"os"
	"path"
	"go_demo/util"
)

type SheetConfig struct {
	// sheet命名
	SheetName string `json:"sheet_name"`
	// 展示数据
	Data [][]interface{} `json:"data"`
}

type ExConfig struct {
	// 路径及文件名
	FileName string `json:"file_name"`
	Sheets []*SheetConfig `json:"sheets"`
}

// 创建excel对象
func CreateExcel(exConfig *ExConfig) (*xlsx.File, error) {
	file := xlsx.NewFile()
	for _, c := range exConfig.Sheets {
		if err := excelAddSheet(c, file); err != nil {
			return nil, err
		}
	}
	return file, nil
}

// 创建excel并本地保存
func CreateExcelToLocal(exConfig *ExConfig) error {
	if exConfig.FileName == "" {
		return errors.New("fileName不能为空")
	}
	file, err := CreateExcel(exConfig)
	if err != nil {
		return err
	}
	if err := file.Save(exConfig.FileName); err != nil {
		return err
	}
	return nil
}

func excelAddSheet(sheetConfig *SheetConfig, file *xlsx.File) error {
	if sheetConfig.SheetName == "" {
		sheetConfig.SheetName = genRandomStr()
	}
	sheet, err := file.AddSheet(sheetConfig.SheetName)
	if err != nil {
		//logger.Error("excel创建sheet报错", "err", err.Error())
		return err
	}
	for _, ds := range sheetConfig.Data {
		row := sheet.AddRow()
		for _, d := range ds {
			row.AddCell().SetValue(d)
		}
	}
	return nil
}

func genRandomStr() string {
	return strconv.Itoa(time.Now().Nanosecond()) + strconv.Itoa(rand.Intn(1000))
}

// 根据二位数组导出excel
func ExportDataToExcel(fPath string, fName string, excelData [][]interface{}) (tmpFilePath string, err error) {
	if !util.FileExist(fPath) {
		if err = os.MkdirAll(fPath, 0755); err != nil {
			return
		}
	}

	tmpFilePath = path.Join(fPath, fName + time.Now().Format("20060102150405") + ".xlsx")
	err = CreateExcelToLocal(&ExConfig{
		FileName: tmpFilePath,
		Sheets: []*SheetConfig {
			{
				// sheet命名
				SheetName: fName,
				// 展示数据
				Data: excelData,
			},
		},
	})
	return
}
