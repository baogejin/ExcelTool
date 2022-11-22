package gencode

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type AnotherExampleCfg struct {
	AnotherSlc []*AnotherInfo `json:"Another"`
	AnotherMap map[int32]*AnotherInfo
}

type AnotherInfo struct {
	ID   int32
	Name string
	Age  int32
}

var anotherexampleCfg *AnotherExampleCfg
var anotherexampleOnce sync.Once

func GetAnotherExampleCfg() *AnotherExampleCfg {
	anotherexampleOnce.Do(func() {
		anotherexampleCfg = new(AnotherExampleCfg)
		anotherexampleCfg.init()
	})
	return anotherexampleCfg
}

func (this *AnotherExampleCfg) init() {
	this.AnotherMap = make(map[int32]*AnotherInfo)
	filePtr, err := os.Open("D:/gs/data/json/AnotherExample.json")
	if err != nil {
		panic(err)
	}
	defer filePtr.Close()
	data, err := ioutil.ReadAll(filePtr)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, this)
	if err != nil {
		panic(err)
	}
	for _, v := range this.AnotherSlc {
		this.AnotherMap[v.ID] = v
	}
}

func (this *AnotherExampleCfg) GetAnotherById(id int32) (*AnotherInfo, bool) {
	if ret, ok := this.AnotherMap[id]; ok {
		return ret, ok
	}
	return nil, false
}
