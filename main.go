package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//go:generate bash ./_run/creator_const_Go.sh
//go:generate bash ./_run/creator_dependencies_Go.sh

// // // // // // // // // // // // // // // // // //

const (
	F1 = "/home/user/GolandProjects/Bookshelf-Writer/TreeOfLanguages/treeLang/en.json"
	F2 = "/home/user/GolandProjects/Bookshelf-Writer/TreeOfLanguages/treeLang/pl.json"
)

func merge(def, data any) any {
	switch defTyped := def.(type) {
	case map[string]any:
		res := make(map[string]any)
		dataMap, ok := data.(map[string]any)
		for k, v := range defTyped {
			if ok {
				if dv, exists := dataMap[k]; exists {
					res[k] = merge(v, dv)
				} else {
					res[k] = merge(v, nil)
				}
			} else {
				res[k] = merge(v, nil)
			}
		}
		return res

	default:
		if data != nil {
			return data
		}
		return ""
	}
}

func main() {
	data1, _ := ioutil.ReadFile(F1)
	var json1 map[string]any
	json.Unmarshal(data1, &json1)
	data2, _ := ioutil.ReadFile(F2)
	var json2 map[string]any
	json.Unmarshal(data2, &json2)

	result := merge(json1, json2)
	out, _ := json.MarshalIndent(result, "", "  ")
	ioutil.WriteFile("file.json", out, os.ModePerm)
}
