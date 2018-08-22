package utils

import (
	"sort"
)

var Array = NewArray()

func NewArray() *array {
	return &array{}
}

type array struct {
}

func (k *array) ChangeKey(handelDatas []map[string]string, key string) (data map[string]interface{}) {
	data = make(map[string]interface{}, len(handelDatas))
	for _, handelData := range handelDatas {
		data[handelData[key]] = handelData
	}
	return
}

func (k *array) ArrayColumn(handelDatas []map[string]string, key string) (data []string) {
	data = make([]string, len(handelDatas))
	for index, handelData := range handelDatas {
		data[index] = handelData[key]
	}
	return
}

func (k *array) ArrayUnique(handelDatas []string) (data []string) {
	sort.Strings(handelDatas)
	for i := 0; i < len(handelDatas); i++ {
		if (i > 0 && handelDatas[i-1] == handelDatas[i]) || len(handelDatas[i]) == 0 {
			continue
		}
		data = append(data, handelDatas[i])
	}
	return
}

func (k *array) InArray(key string, search []string) bool {
	if len(search) == 0 {
		return false
	}
	var searchMap = map[string]bool{}
	for _, s := range search {
		searchMap[s] = true
	}
	_, ok := searchMap[key]
	return ok
}