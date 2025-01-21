package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// // // // // // // // // // // // // // // // // //

func _deepCopy(src any) map[string]any {
	tempMap := make(map[string]any)
	bytes, _ := json.Marshal(src)
	json.Unmarshal(bytes, &tempMap)
	return tempMap
}

func sortedJSON(data any) []byte {
	tempMap := _deepCopy(data)

	keys := make([]string, 0, len(tempMap))
	for key := range tempMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sortedMap := make(map[string]any)
	for _, key := range keys {
		sortedMap[key] = tempMap[key]
	}

	retData, _ := json.MarshalIndent(sortedMap, "", "  ")
	return retData
}

func saveTempFileJSON(data any) string {
	json := sortedJSON(data)

	tempFile, _ := os.CreateTemp(os.TempDir(), fmt.Sprintf("%d_json_temp", time.Now().Nanosecond()))
	defer tempFile.Close()

	tempFile.Write(json)

	return tempFile.Name()
}
