package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
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

func diffLines(file1, file2 string) (map[int]string, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return nil, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return nil, err
	}
	defer f2.Close()
	differences := make(map[int]string)

	scanner1 := bufio.NewScanner(f1)
	scanner2 := bufio.NewScanner(f2)

	lineNumber := 1
	for scanner1.Scan() && scanner2.Scan() {
		line1 := scanner1.Text()
		line2 := scanner2.Text()

		if line1 != line2 {
			differences[lineNumber] = strings.TrimSpace(line2)
		}
		lineNumber++
	}

	if err := scanner1.Err(); err != nil {
		return nil, err
	}
	if err := scanner2.Err(); err != nil {
		return nil, err
	}

	for scanner2.Scan() {
		differences[lineNumber] = strings.TrimSpace(scanner2.Text())
		lineNumber++
	}

	return differences, nil
}

func saveSortedJSON(data any, filename string) error {
	tempMap := make(map[string]any)

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
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

	sortedJSON, err := json.MarshalIndent(sortedMap, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(sortedJSON)
	if err != nil {
		return err
	}

	return nil
}

// //

// Функция очищает структуру, устанавливая полям значения по умолчанию
func clearStruct(data any) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr {
		panic("clearStruct requires a pointer")
	}

	elem := val.Elem()
	clearValue(elem)
}

// Функция рекурсивно очищает значения
func clearValue(value reflect.Value) {
	switch value.Kind() {

	case reflect.Ptr:
		if !value.IsNil() {
			clearValue(value.Elem())
		}

	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.CanSet() {
				clearValue(field)
			}
		}

	case reflect.Interface:
		if !value.IsNil() {
			elem := value.Elem()
			clearValue(elem)
			// value.Set(reflect.Zero(value.Type()))
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			bufVal := value.MapIndex(key)

			_, ok := bufVal.Interface().(string)
			if ok {
				value.SetMapIndex(key, reflect.Zero(reflect.TypeOf("")))
			} else {
				clearValue(bufVal)
			}
		}

	default:
		if !value.IsValid() || !value.CanSet() {
			return
		}
		value.Set(reflect.Zero(value.Type()))
	}
}

// //

func main() {
	data1, _ := ioutil.ReadFile(F1)
	var json1 map[string]any
	json.Unmarshal(data1, &json1)
	data2, _ := ioutil.ReadFile(F2)
	var json2 map[string]any
	json.Unmarshal(data2, &json2)

	result := merge(json1, json2)

	saveSortedJSON(result, "file.json")

	//

	ff1 := "file.json"
	ff2 := "file_f.json"

	clearStruct(&json1)
	saveSortedJSON(json1, ff1)

	clearStruct(&json2)
	saveSortedJSON(json2, ff2)

	rf1, _ := os.ReadFile(ff1)
	rf2, _ := os.ReadFile(ff2)

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(rf1)),
		B:        difflib.SplitLines(string(rf2)),
		FromFile: ff1,
		ToFile:   ff2,
		Context:  2,
	}

	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		fmt.Println("Ошибка при генерации диффа:", err)
		return
	}

	fmt.Println(diffText)
}
