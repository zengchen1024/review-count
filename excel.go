package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

const sheetName = "Sheet1"

func newExcel(fileName string) (*excel, error) {
	if _, err := os.Stat(fileName); err == nil || !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s exists or internal error", fileName)
	}

	v := &excel{
		file:     excelize.NewFile(),
		row:      1,
		fileName: fileName,
	}

	err := v.write(&reviewRecord{
		pr:         "PR",
		commentURL: "comment URL",
		comment:    "comment",
	})
	if err != nil {
		if err1 := v.file.Close(); err1 != nil {
			fmt.Printf("clean excel failed, err:%v\n", err)
		}
	}

	return v, err
}

type excel struct {
	file     *excelize.File
	row      int
	fileName string
}

func (w *excel) write(record *reviewRecord) (err error) {
	row := strconv.Itoa(w.row)
	w.row++

	err = w.file.SetCellValue(sheetName, "A"+row, record.pr)
	err = w.file.SetCellValue(sheetName, "B"+row, record.commentURL)
	err = w.file.SetCellValue(sheetName, "C"+row, record.comment)

	return
}

func (w *excel) save() error {
	if err := w.file.SaveAs(w.fileName); err != nil {
		return err
	}

	return w.file.Close()
}
