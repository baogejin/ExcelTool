package myexcel

import (
	"errors"
	"strconv"
	"strings"
)

type CellType int32

const (
	CellTypeSimple    CellType = 0
	CellTypeSlc       CellType = 1
	CellTypeDoubleSlc CellType = 2
	CellTypeMap       CellType = 3
)

type TypeInfo struct {
	CType      CellType
	ValueType1 string
	ValueType2 string
}

func (this *TypeInfo) FixType() {
	if this.ValueType1 == "int" {
		this.ValueType1 = "int32"
	} else if this.ValueType1 == "float" {
		this.ValueType1 = "float32"
	}
	if this.ValueType2 == "int" {
		this.ValueType2 = "int32"
	} else if this.ValueType2 == "float" {
		this.ValueType2 = "float32"
	}
}

func (this *TypeInfo) ParseToJson(str string) (string, error) {
	switch this.CType {
	case CellTypeSimple:
		return parseValueStr(this.ValueType1, str)
	case CellTypeSlc:
		if str == "" {
			return "[]", nil
		}
		ret := "["
		s := strings.Split(str, "|")
		for i, v := range s {
			if i != 0 {
				ret += ","
			}
			if vStr, err := parseValueStr(this.ValueType1, v); err == nil {
				ret += vStr
			} else {
				return "", err
			}
		}
		ret += "]"
		return ret, nil
	case CellTypeDoubleSlc:
		if str == "" {
			return "[[]]", nil
		}
		ret := "["
		s1 := strings.Split(str, "|")
		for i, s := range s1 {
			if i != 0 {
				ret += ","
			}
			if s == "" {
				ret += "[]"
				continue
			}
			ret += "["
			s2 := strings.Split(s, ":")
			for j, v := range s2 {
				if j != 0 {
					ret += ","
				}
				if vStr, err := parseValueStr(this.ValueType1, v); err == nil {
					ret += vStr
				} else {
					return "", err
				}
			}
			ret += "]"
		}
		ret += "]"
		return ret, nil
	case CellTypeMap:
		if str == "" {
			return "{}", nil
		}
		ret := "{"
		s1 := strings.Split(str, "|")
		repeatCheck := make(map[string]bool)
		for i, s := range s1 {
			if i != 0 {
				ret += ","
			}
			if s == "" {
				continue
			}
			s2 := strings.Split(s, ":")

			if len(s2) != 2 {
				return "", errors.New("map type len err")
			}
			if repeatCheck[s2[0]] {
				return "", errors.New("map type key repeat")
			}
			repeatCheck[s2[0]] = true
			if _, err := parseValueStr(this.ValueType1, s2[0]); err == nil {
				ret += "\"" + s2[0] + "\""
			} else {
				return "", err
			}
			ret += ":"
			if vStr, err := parseValueStr(this.ValueType2, s2[1]); err == nil {
				ret += vStr
			} else {
				return "", err
			}
		}
		ret += "}"
		return ret, nil
	}
	return "", errors.New("cell type err")
}

func parseValueStr(vType, value string) (string, error) {
	switch vType {
	case "bool":
		value = strings.ToLower(value)
		if value == "true" || value == "false" {
			return value, nil
		}
		if num, err := strconv.ParseInt(value, 10, 0); err == nil {
			if num == 0 {
				return "false", nil
			} else {
				return "true", nil
			}
		} else {
			return "false", nil
		}
	case "int32", "int64":
		if _, err := strconv.ParseInt(value, 10, 0); err == nil {
			return value, nil
		} else {
			return "0", nil
		}
	case "float32", "float64":
		if _, err := strconv.ParseFloat(value, 64); err == nil {
			return value, nil
		} else {
			return "0", nil
		}
	case "string":
		return "\"" + value + "\"", nil
	}
	return "", errors.New("value type invalid:" + vType)
}

func getTypeInfoByStr(str string) (*TypeInfo, error) {
	err := errors.New("type info invalid " + str)
	l := len(str)
	if checkValueTypeValid(str) {
		return &TypeInfo{
			CType:      CellTypeSimple,
			ValueType1: str,
		}, nil
	} else if strings.HasSuffix(str, "[][]") {
		if l < 5 || !checkValueTypeValid(str[:l-4]) {
			return nil, err
		}
		return &TypeInfo{
			CType:      CellTypeDoubleSlc,
			ValueType1: str[:l-4],
		}, nil
	} else if strings.HasPrefix(str, "double_slc|") {
		if l < 12 || !checkValueTypeValid(str[11:]) {
			return nil, err
		}
		return &TypeInfo{
			CType:      CellTypeDoubleSlc,
			ValueType1: str[11:],
		}, nil
	} else if strings.HasSuffix(str, "[]") {
		if l < 3 || !checkValueTypeValid(str[:l-2]) {
			return nil, err
		}
		return &TypeInfo{
			CType:      CellTypeSlc,
			ValueType1: str[:l-2],
		}, nil
	} else if strings.HasPrefix(str, "slc|") {
		if l < 5 || !checkValueTypeValid(str[4:]) {
			return nil, err
		}
		return &TypeInfo{
			CType:      CellTypeSlc,
			ValueType1: str[4:],
		}, nil
	} else if strings.HasPrefix(str, "map|") {
		splits := strings.Split(str, "|")
		if len(splits) != 3 {
			return nil, err
		}
		if !checkValueTypeValid(splits[1]) || !checkValueTypeValid(splits[2]) {
			return nil, err
		}
		return &TypeInfo{
			CType:      CellTypeMap,
			ValueType1: splits[1],
			ValueType2: splits[2],
		}, nil
	}
	return nil, err
}

func checkValueTypeValid(typeStr string) bool {
	switch typeStr {
	case "bool", "int", "int32", "int64", "float", "float32", "float64", "string":
		return true
	default:
		return false
	}
}
