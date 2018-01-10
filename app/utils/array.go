package utils

import (
	"sort"
)

func NewArray() *Array {
	return &Array{}
}

type Array struct {
}

func (k *Array) ChangeKey(handelDatas []map[string]string, key string) (data map[string]interface{}) {
	data = make(map[string]interface{}, len(handelDatas))
	for _, handelData := range handelDatas {
		data[handelData[key]] = handelData
	}
	return
}

func (k *Array) ArrayColumn(handelDatas []map[string]string, key string) (data []string) {
	data = make([]string, len(handelDatas))
	for index, handelData := range handelDatas {
		data[index] = handelData[key]
	}
	return
}

func (k *Array) ArrayUnique(handelDatas []string) (data []string) {
	sort.Strings(handelDatas)
	for i := 0; i < len(handelDatas); i++ {
		if (i > 0 && handelDatas[i-1] == handelDatas[i]) || len(handelDatas[i]) == 0 {
			continue
		}
		data = append(data, handelDatas[i])
	}
	return
}
