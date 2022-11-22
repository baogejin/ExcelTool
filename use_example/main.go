package main

import (
	"ExcelTool/data/gencode"
	"fmt"
)

func main() {
	if e, ok := gencode.GetExampleCfg().GetExampleById(1); ok {
		fmt.Println(e.Name)
	}
	if a, ok := gencode.GetAnotherExampleCfg().GetAnotherById(1); ok {
		fmt.Println(a.Name, a.Age)
	}
}
