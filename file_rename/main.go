package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/go_multi_log/logger/levels"
	"github.com/takecontrolsoft/go_multi_log/logger/loggers"
	"github.com/xuri/excelize/v2"
)

func main() {
	fl, shouldReturn := RegisterLogger()
	if shouldReturn {
		return
	}
	defer func() {
		fl.Stop()
	}()

	m, shouldReturn := ReadExcel("Clients.xlsx", "Clients", "ЕИК", "ИМЕ")
	if shouldReturn {
		return
	}
	currentDir, err := filepath.Abs(".")
	if err != nil {
		logger.Error("Getting current directory failed.")
		return
	}
	subitems, err := os.ReadDir(currentDir)
	if err != nil {
		logger.ErrorF("Getting files from %s failed", currentDir)
		return
	}
	for _, subitem := range subitems {
		fn := subitem.Name()
		var newName = fn
		fileExtension := filepath.Ext(fn)
		if !subitem.IsDir() && strings.ToLower(fileExtension) == ".txt" {
			logger.InfoF("File: %s", fn)
			for eik, name := range m {
				if strings.Contains(fn, eik) && !strings.Contains(strings.ToUpper(fn), strings.ToUpper(name)) {
					var dekl = 0
					if strings.Contains(strings.ToUpper(newName), "EMPL2021") {
						dekl = 1
						newName = strings.ReplaceAll(newName, "EMPL2021", "")
						newName = strings.ReplaceAll(newName, "empl2021", "")
					} else if strings.Contains(strings.ToUpper(newName), "NRA62007") {
						dekl = 6
						newName = strings.ReplaceAll(newName, "NRA62007", "")
						newName = strings.ReplaceAll(newName, "nra62007", "")
					}
					if dekl == 0 {
						newName = strings.ReplaceAll(newName, eik, fmt.Sprintf("%s_%s", name, eik))
					} else {
						newName = strings.TrimLeft(newName, "_")
						newName = strings.ReplaceAll(newName, eik, fmt.Sprintf("%s_(дек. %v)_%s", name, dekl, eik))
					}
					os.Rename(filepath.Join(currentDir, fn), filepath.Join(currentDir, newName))
					logger.InfoF("File '%s' renamed to '%s' in folder %s", fn, newName, currentDir)
				}

			}
		}
	}
	logger.InfoF("File renaming completed successfully")

}

func ReadExcel(excelName string, sheetName string, eik string, name string) (map[string]string, bool) {
	var m = make(map[string]string)
	fe, err := excelize.OpenFile(excelName)
	if err != nil {
		logger.ErrorF("Open file %s failed, Error: %v", excelName, err)
		return nil, true
	}
	defer func() {

		if err := fe.Close(); err != nil {
			logger.ErrorF("Closing file %s failed, Error: %v", excelName, err)
		}
	}()

	rows, err := fe.GetRows(sheetName)
	if err != nil {
		logger.ErrorF("Getting rows from %s failed, Error: %v", sheetName, err)
		return nil, true
	}
	var eikIndex, nameIndex int = 0, 0

	for c, colCell := range rows[0] {
		logger.InfoF("Row headers: %v", colCell)
		if colCell == eik {
			eikIndex = c
		} else if colCell == name {
			nameIndex = c
		}
	}
	logger.InfoF("EIK index: %v, NAME index: %v", eikIndex, nameIndex)

	for r, row := range rows {
		if r > 0 {
			logger.InfoF("Row index: %v", r)
			eikValue := row[eikIndex]
			nameValue := row[nameIndex]

			m[eikValue] = nameValue
			logger.InfoF("EIK: %s, NAME: %s", eikValue, nameValue)
		}
	}
	return m, false
}

func RegisterLogger() (*loggers.FileLogger, bool) {
	fileOptions := loggers.FileOptions{
		Directory:     "./logs",
		FilePrefix:    "file_rename",
		FileExtension: ".log",
	}
	err := os.MkdirAll(fileOptions.Directory, os.ModePerm)
	if err != nil {
		return nil, false
	}
	fl := loggers.NewFileLogger(levels.All, "", fileOptions)

	err = logger.RegisterLogger("file_rename", fl)
	if err != nil {
		logger.ErrorF("Creating log file failed %v", err)
		return nil, true
	}
	return fl, false
}
