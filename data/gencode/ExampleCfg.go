package gencode

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type ExampleCfg struct {
	ExampleSlc []*ExampleInfo `json:"Example"`
	ExampleMap map[int32]*ExampleInfo
	AbbSlc     []*AbbInfo `json:"Abb"`
	AbbMap     map[int32]*AbbInfo
}

type ExampleInfo struct {
	ID         int32
	Name       string
	Slc1       []int32
	Slc2       []float32
	DoubleSlc1 [][]int32
	DoubleSlc2 [][]string
	IsBool     bool
	Map1       map[int32]int32
	Map2       map[int32]string
}

type AbbInfo struct {
	ID  int32
	Sth string
}

var exampleCfg *ExampleCfg
var exampleOnce sync.Once

func GetExampleCfg() *ExampleCfg {
	exampleOnce.Do(func() {
		exampleCfg = new(ExampleCfg)
		exampleCfg.init()
	})
	return exampleCfg
}

func (this *ExampleCfg) init() {
	this.ExampleMap = make(map[int32]*ExampleInfo)
	this.AbbMap = make(map[int32]*AbbInfo)
	filePtr, err := os.Open("D:/gs/data/json/Example.json")
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
	for _, v := range this.ExampleSlc {
		this.ExampleMap[v.ID] = v
	}
	for _, v := range this.AbbSlc {
		this.AbbMap[v.ID] = v
	}
}

func (this *ExampleCfg) GetExampleById(id int32) (*ExampleInfo, bool) {
	if ret, ok := this.ExampleMap[id]; ok {
		return ret, ok
	}
	return nil, false
}

func (this *ExampleCfg) GetAbbById(id int32) (*AbbInfo, bool) {
	if ret, ok := this.AbbMap[id]; ok {
		return ret, ok
	}
	return nil, false
}
