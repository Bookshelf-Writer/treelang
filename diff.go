package main

import (
	"encoding/json"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
)

// // // // // // // // // // // // // // // // // //

func _saveSortedJSON(data any) string {
	tempMap := make(map[string]any)

	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &tempMap)

	keys := make([]string, 0, len(tempMap))
	for key := range tempMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sortedMap := make(map[string]any)
	for _, key := range keys {
		sortedMap[key] = tempMap[key]
	}

	sortedJSON, _ := json.MarshalIndent(sortedMap, "", "  ")

	tempFile, _ := os.CreateTemp(os.TempDir(), fmt.Sprintf("%d_json_temp", time.Now().Nanosecond()))
	defer tempFile.Close()

	tempFile.Write(sortedJSON)

	return tempFile.Name()
}

func deepCopy(src any) any {
	srcVal := reflect.ValueOf(src)
	copiedVal := reflect.New(srcVal.Type()).Elem()
	copiedVal.Set(srcVal)
	return copiedVal.Interface()
}

func diff(def, data any) ([]string, error) {
	defValue := deepCopy(def)
	dataValue := deepCopy(data)

	clearObj(&defValue)
	file1 := _saveSortedJSON(defValue)

	clearObj(&dataValue)
	file2 := _saveSortedJSON(dataValue)

	rf1, _ := os.ReadFile(file1)
	rf2, _ := os.ReadFile(file2)

	dd := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(rf2)),
		B:        difflib.SplitLines(string(rf1)),
		FromFile: file1,
		ToFile:   file2,
		Context:  1,
	}
	diffText, err := difflib.GetUnifiedDiffString(dd)
	if err != nil {
		return nil, err
	}

	bufArr := strings.Split(diffText, "@@")
	finalArr := make([]string, 0, len(bufArr)/2)
	bufArr = bufArr[1:]

	for x := 0; x < len(bufArr); x = x + 2 {
		txt := bufArr[x+1]
		txt = strings.Replace(txt, "\"\"", "...", -1)
		finalArr = append(finalArr, txt)
	}

	return finalArr, nil
}
