package main

import (
	"ExcelTool/myexcel"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func main() {
	excelPath := "D:/ExcelTool/data/excel"
	jsonPath := "D:/ExcelTool/data/json"
	codePath := "D:/ExcelTool/data/gencode"
	fileList, err := getExcelList(excelPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	//读取excel
	excels := loadExcels(excelPath, fileList)
	//生成json
	genJson(jsonPath, excels)
	//生成go代码
	genGoCode(codePath, excels)

}

func getExcelList(path string) ([]string, error) {
	ret := []string{}
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		return ret, err
	}
	for _, fileInfo := range fileInfoList {
		name := fileInfo.Name()
		if strings.HasPrefix(name, "~") {
			continue
		}
		strSlc := strings.Split(name, ".")
		if len(strSlc) != 2 || strSlc[1] != "xlsx" {
			continue
		}
		if len([]byte(name)) != len([]rune(name)) {
			fmt.Println("文件名请不要用英文数字以外的字符:", name)
			continue
		}
		ret = append(ret, strSlc[0])
	}
	return ret, nil
}

func loadExcels(path string, files []string) []*myexcel.ExcelInfo {
	ret := []*myexcel.ExcelInfo{}
	for _, file := range files {
		excel := &myexcel.ExcelInfo{}
		excel.Name = file
		err := excel.Load(path, file)
		if err != nil {
			fmt.Println("load ", file, ".xlsx failed,", err)
			continue
		}
		ret = append(ret, excel)
		fmt.Println("load", file+".xlsx success")
	}
	return ret
}

func genJson(path string, excels []*myexcel.ExcelInfo) {
	for _, excel := range excels {
		if err := excel.GenJson(path); err != nil {
			fmt.Println("gen json failed:"+excel.Name+".xlsx", err)
		}
	}
}

func genGoCode(path string, excels []*myexcel.ExcelInfo) {
	for _, excel := range excels {
		if err := excel.GenCode(path); err != nil {
			fmt.Println("gen code failed:"+excel.Name+".xlsx", err)
		}
	}
	cmd := exec.Command("gofmt", "-w", path)
	err := cmd.Run()
	if nil != err {
		fmt.Println(err)
	}
}
