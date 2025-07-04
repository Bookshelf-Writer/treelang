package generate

import (
	"fmt"
	. "github.com/Bookshelf-Writer/treelang"
	"strings"
)

// // // // // // // // // // // // // // // // // //

func merge(def, data any, pad int, key string) any {
	printMsg := func(format string, a ...any) {
		fmt.Printf(strings.Repeat("\t", pad)+format+"\n", a...)
	}

	if defMap, ok := def.(map[string]any); ok {
		result := make(map[string]any)

		dataMap, dataIsMap := data.(map[string]any)
		for k, defVal := range defMap {
			kk := strings.Join([]string{key, k}, ".")

			if dataIsMap {
				if dataVal, exists := dataMap[k]; exists {
					result[k] = merge(defVal, dataVal, pad, kk)
				} else {
					printMsg("%s: %s", Cyan(kk), Red("no field in file"))
					result[k] = merge(defVal, nil, pad, kk)
				}
			} else {
				result[k] = merge(defVal, nil, pad, kk)
			}
		}
		return result
	}

	if _, ok := def.(string); ok {
		if data == nil {
			return key
		}

		if dataStr, ok := data.(string); ok {
			return dataStr
		}
		printMsg("%s: %s", Cyan(key), Yellow("must be a string"))
		return key
	}

	if _, ok := def.([]any); ok {
		if dataStr, ok := data.([]any); ok {
			if len(dataStr) > 0 {
				if _, ok := dataStr[0].(string); ok {
					return dataStr
				}
				printMsg("%s: %s", Cyan(key), Yellow("the array must be made up of strings"))
			}
			return []string{}
		}
		printMsg("%s: %s", Cyan(key), Yellow("must be an array"))
		return []string{}
	}

	return key
}

func mergeLangObj(master, slave *LangObj, pad int) *LangObj {
	data := merge(master.Data, slave.Data, pad, "data")
	slave.Data = data.(map[string]any)
	return slave
}
