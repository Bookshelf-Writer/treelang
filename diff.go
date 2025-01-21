package main

import (
	"github.com/pmezard/go-difflib/difflib"
	"os"
	"strings"
)

// // // // // // // // // // // // // // // // // //

func diffFile(filepathA, filepathB string) ([]string, error) {
	rf1, err := os.ReadFile(filepathA)
	if err != nil {
		return nil, err
	}

	rf2, err := os.ReadFile(filepathB)
	if err != nil {
		return nil, err
	}

	dd := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(rf2)),
		B:        difflib.SplitLines(string(rf1)),
		FromFile: filepathA,
		ToFile:   filepathB,
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
		finalArr = append(finalArr, bufArr[x+1])
	}

	return finalArr, nil
}

func diffObj(def, data any) ([]string, error) {
	defValue := _deepCopy(def)
	dataValue := _deepCopy(data)

	clearObj(&defValue)
	file1 := saveTempFileJSON(defValue)

	clearObj(&dataValue)
	file2 := saveTempFileJSON(dataValue)

	arr, err := diffFile(file1, file2)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(arr); i++ {
		txt := arr[i]
		txt = strings.Replace(txt, "\"\"", "...", -1)
		arr[i] = txt
	}

	return arr, nil
}
