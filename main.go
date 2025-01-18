package main

import (
	"encoding/json"
	"fmt"
	"os"
)

//go:generate bash ./_run/creator_const_Go.sh
//go:generate bash ./_run/creator_dependencies_Go.sh

// // // // // // // // // // // // // // // // // //

const (
	F1 = "treeLang/en.json"
	F2 = "treeLang/pl.json"
)

// //

func main() {
	data1, _ := os.ReadFile(F1)
	var json1 map[string]any
	json.Unmarshal(data1, &json1)
	data2, _ := os.ReadFile(F2)
	var json2 map[string]any
	json.Unmarshal(data2, &json2)

	result := merge(json1, json2)
	fmt.Println(result)

	//

	arr, err := diffObj(json1, json2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(arr))
	for _, txt := range arr {
		fmt.Println("\n", txt)
	}
}
