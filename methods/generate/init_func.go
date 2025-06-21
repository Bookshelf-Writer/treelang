package generate

import (
	"fmt"
	. "github.com/Bookshelf-Writer/treelang"
	"os"
	"path/filepath"
)

// // // // // // // // // // // // // // // // // //

func parseSlave(fromFilePath, fromReadDir string) (*LangObj, map[string]*LangObj, error) {
	master, err := ReadFile(fromFilePath)
	if err != nil {
		return nil, nil, err
	}

	files, err := os.ReadDir(fromReadDir)
	if err != nil {
		return master, nil, err
	}

	slave := make(map[string]*LangObj, 0)
	for _, file := range files {
		if !file.IsDir() {
			obj, err := ReadFile(filepath.Join(fromReadDir, file.Name()))
			if err != nil {
				continue
			}

			if obj.Info != nil && obj.Data != nil {
				slave[file.Name()] = obj
			}
		}
	}

	return master, slave, nil
}

func parseMap(fromReadDir string) ([]*LangInfoObj, error) {
	files, err := os.ReadDir(fromReadDir)
	if err != nil {
		return nil, err
	}

	arr := make([]*LangInfoObj, 0)
	for _, file := range files {
		if !file.IsDir() {
			obj, err := ReadFile(filepath.Join(fromReadDir, file.Name()))
			if err != nil {
				continue
			}

			if obj.Info != nil && obj.Data != nil {
				arr = append(arr, obj.Info)
			}
		}
	}

	return arr, nil
}

// // // //

func writeJsonData(fromFilePath, fromReadDir, toDir string) error {
	master, slave, err := parseSlave(fromFilePath, fromReadDir)
	if err != nil {
		return err
	}

	for fileName, obj := range slave {
		fmt.Printf("%s:\n", Green(fileName))
		finalObj := mergeLangObj(master, obj, 1)
		err = createLangJSON(finalObj, toDir)
		if err == nil {
			fmt.Printf("Created: JSON %s\n", Blue(finalObj.Info.Name.EN))
		}
	}

	return nil
}

func writeJsonMap(fromReadDir, toDir string) error {
	arr, err := parseMap(fromReadDir)
	if err != nil {
		return err
	}

	err = createMapJSON(arr, toDir)
	if err == nil {
		fmt.Printf("Created: MAP\n")
	}

	return err
}

// //

func writeYmlData(fromFilePath, fromReadDir, toDir string) error {
	master, slave, err := parseSlave(fromFilePath, fromReadDir)
	if err != nil {
		return err
	}

	for fileName, obj := range slave {
		fmt.Printf("%s:\n", Green(fileName))
		finalObj := mergeLangObj(master, obj, 1)
		err = createLangYML(finalObj, toDir)
		if err == nil {
			fmt.Printf("Created: YML %s\n", Blue(finalObj.Info.Name.EN))
		}
	}

	return nil
}

func writeYmlMap(fromReadDir, toDir string) error {
	arr, err := parseMap(fromReadDir)
	if err != nil {
		return err
	}

	err = createMapYML(arr, toDir)
	if err == nil {
		fmt.Printf("Created: MAP\n")
	}

	return err
}

// //

func writeGoData(fromFilePath, fromReadDir, toDir, packageName string) error {
	master, slave, err := parseSlave(fromFilePath, fromReadDir)
	if err != nil {
		return err
	}

	for fileName, obj := range slave {
		fmt.Printf("%s:\n", Green(fileName))
		finalObj := mergeLangObj(master, obj, 1)
		err = createLangGO(finalObj, toDir, packageName)
		if err == nil {
			fmt.Printf("Created: GO %s\n", Blue(finalObj.Info.Name.EN))
		}
	}

	return writeGoStruct(fromFilePath, toDir, packageName)
}

func writeGoMap(fromReadDir, toDir, packageName string) error {
	arr, err := parseMap(fromReadDir)
	if err != nil {
		return err
	}

	err = createMapGO(arr, toDir, packageName)
	if err == nil {
		fmt.Printf("Created: MAP\n")
	}

	return err
}
