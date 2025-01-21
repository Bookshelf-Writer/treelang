package main

import (
	"errors"
	"fmt"
	generator "github.com/Bookshelf-Writer/SimpleGenerator"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

// // // // // // // // // // // // // // // // // //

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generating files from a tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		var fromPath, masterPath, toPath FilePathType
		var err error

		if CmdFromFilePath != "" {
			fromPath, err = CheckPathCMD(&CmdFromFilePath, "from")
			if err != nil {
				return err
			}
		}
		if CmdToFilePath != "" {
			toPath, err = CheckPathCMD(&CmdToFilePath, "to")
			if err != nil {
				return err
			}
			if toPath != FilePathIsDir {
				return paramErr("to", errors.New("should only point to a directory"))
			}
		}
		if CmdMasterFile != "" {
			masterPath, err = CheckPathCMD(&CmdMasterFile, "master")
			if err != nil {
				return err
			}
			if masterPath != FilePathValid {
				return paramErr("master", errors.New("the path should only point to an existing file"))
			}
		}

		if fromPath == FilePathUnknown {
			return paramErr("from", errors.New("this parameter is required"))
		}
		if toPath == FilePathUnknown {
			return paramErr("to", errors.New("this parameter is required"))
		}

		if fromPath == FilePathIsDir {
			if masterPath == FilePathUnknown {
				return paramErr("master", errors.New("this parameter is required"))
			}
		}
		if fromPath != FilePathValid && fromPath != FilePathIsDir {
			return paramErr("from", errors.New("error with parameter. must point to a file or folder"))
		}
		if fromPath == FilePathValid {
			if CmdMasterFile == "" {
				CmdMasterFile = CmdFromFilePath
				masterPath = FilePathValid
			}
			CmdFromFilePath = path.Dir(CmdFromFilePath)
		}
		if masterPath != FilePathValid {
			return paramErr("master", errors.New("error with parameter. must point to a file"))
		}

		// //

		errFile := errors.New("the specified master-file is not a valid language file")
		obj, err := ReadFile(CmdMasterFile)
		if err != nil {
			return errFile
		}
		if obj == nil {
			return errFile
		}
		if obj.Data == nil || obj.Info == nil {
			return errFile
		}
		if obj.Info.Name == nil {
			return errFile
		}

		// //

		if !*CmdJson && !*CmdYml && CmdPackage == "" {
			return fmt.Errorf("generation type is not specified. Select %s, %s or enter %s", cyan("--yml"), cyan("--json"), cyan("----go-package"))
		}

		// //

		if *CmdJson {
			if *CmdMap {
				return writeJsonMap(CmdFromFilePath, CmdToFilePath)
			} else {
				return writeJsonData(CmdMasterFile, CmdFromFilePath, CmdToFilePath)
			}
		}

		if *CmdYml {
			if *CmdMap {
				return writeYmlMap(CmdFromFilePath, CmdToFilePath)
			} else {
				return writeYmlData(CmdMasterFile, CmdFromFilePath, CmdToFilePath)
			}
		}

		if CmdPackage != "" {
			if *CmdMap {
				return writeGoMap(CmdFromFilePath, CmdToFilePath, CmdPackage)
			} else {
				return writeGoData(CmdMasterFile, CmdFromFilePath, CmdToFilePath, CmdPackage)
			}
		}

		return errors.New("Unexpected termination of context")
	},
}

// // // // // //

func writeJsonData(fromFilePath, fromReadDir, toDir string) error {
	master, err := ReadFile(fromFilePath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(fromReadDir)
	if err != nil {
		return err
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

	for fileName, obj := range slave {
		fmt.Printf("%s:\n", green(fileName))
		finalObj := mergeLangObj(master, obj, 1)
		err = createLangJSON(finalObj, toDir)
		if err == nil {
			fmt.Printf("Created: JSON %s\n", blue(finalObj.Info.Name.EN))
		}
	}

	return nil
}

func writeJsonMap(fromReadDir, toDir string) error {
	files, err := os.ReadDir(fromReadDir)
	if err != nil {
		return err
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

	err = createMapJSON(arr, toDir)
	if err == nil {
		fmt.Printf("Created: MAP\n")
	}

	return err
}

// //

func writeYmlData(fromFilePath, fromReadDir, toDir string) error {
	master, err := ReadFile(fromFilePath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(fromReadDir)
	if err != nil {
		return err
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

	for fileName, obj := range slave {
		fmt.Printf("%s:\n", green(fileName))
		finalObj := mergeLangObj(master, obj, 1)
		err = createLangYML(finalObj, toDir)
		if err == nil {
			fmt.Printf("Created: YML %s\n", blue(finalObj.Info.Name.EN))
		}
	}

	return nil
}

func writeYmlMap(fromReadDir, toDir string) error {
	files, err := os.ReadDir(fromReadDir)
	if err != nil {
		return err
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

	err = createMapYML(arr, toDir)
	if err == nil {
		fmt.Printf("Created: MAP\n")
	}

	return err
}

// //

func prepareStructMap(m map[string]any) map[string]string {
	result := make(map[string]string)

	for key, value := range m {
		switch v := value.(type) {
		case map[string]any:
			result[key] = "obj"
			nestedMap := prepareStructMap(v)
			for nestedKey, nestedType := range nestedMap {
				result[key+":"+nestedKey] = nestedType
			}
		default:
			switch reflect.TypeOf(v).String() {
			case "string":
				result[key] = "str"
			case "[]interface {}":
				result[key] = "arr"
			default:
				result[key] = reflect.TypeOf(v).String()
			}
		}
	}

	return result
}

func structMap(data map[string]any) map[string]map[string][]string {
	dataMap := prepareStructMap(data)
	kkk := make(map[string]map[string][]string)

	baseDef := ToGoVariableName("Data")

	for fullKey, val := range dataMap {
		keys := strings.Split(fullKey, ":")

		var base, field, name string
		if len(keys) == 1 {
			base = baseDef
			name = keys[0]
			field = ToGoVariableName(keys[0])
		} else {
			base = ToGoVariableName(keys[len(keys)-2])
			name = keys[len(keys)-1]
			field = ToGoVariableName(keys[len(keys)-1])
		}
		if _, ok := kkk[base]; !ok {
			kkk[base] = make(map[string][]string)
		}

		switch val {
		case "str":
			kkk[base][field] = []string{"string", name}
		case "arr":
			kkk[base][field] = []string{"[]string", name}
		case "obj":
			kkk[base][field] = []string{"*" + field + "Obj", name}
		}
	}

	return kkk
}

func writeGoStruct(fromFilePath, toDir, packageName string) error {
	goGen := generator.NewFilePathName(toDir, packageName)
	types := func(name, jsonName string) generator.GeneratorTypeObj {
		return generator.GeneratorTypeObj{Types: goGen.NewType(name), Tags: map[string]string{"json": jsonName, "yaml": jsonName}}
	}

	goGen.NewImport("bytes", "")
	goGen.NewImport("github.com/srwiley/oksvg", "")
	goGen.NewImport("github.com/srwiley/rasterx", "")
	goGen.NewImport("image", "")

	// //

	obj, err := ReadFile(fromFilePath)
	if err != nil {
		return err
	}

	kkk := structMap(obj.Data)

	for _, key := range sortMapKey(kkk) {
		bufMap := make(map[string]generator.GeneratorTypeObj)

		for _, k := range sortMapKey(kkk[key]) {
			bufMap[k] = types(kkk[key][k][0], kkk[key][k][1])
		}
		goGen.AddStruct(key, bufMap)
	}

	goGen.SeparatorX4().LN()

	// //

	goGen.AddStruct("LangInfoName", map[string]generator.GeneratorTypeObj{
		"EN":  types("string", "en"),
		"DEF": types("string", "def"),
	})

	infoObj := goGen.AddStruct("LangInfo", map[string]generator.GeneratorTypeObj{
		"Name": types("*LangInfoNameObj", "name"),
		"Code": types("string", "code"),
		"Flag": types("string", "flag"),
	})

	goGen.AddStruct("LangSys", map[string]generator.GeneratorTypeObj{
		"Date": types("string", "date"),
		"Hash": types("string", "hash"),
	})

	goGen.AddFunc(
		"PNG",
		nil,
		map[string]generator.GeneratorTypeObj{
			"err": generator.GeneratorTypeObj{Types: goGen.TypeError()},
			"img": generator.GeneratorTypeObj{Types: goGen.NewType("*image.RGBA")},
		},
		infoObj,
		func(gen *generator.GeneratorObj) {
			gen.PrintLN("svgReader := bytes.NewReader([]byte(parent.Flag))")
			gen.PrintLN("icon, err := oksvg.ReadIconStream(svgReader)")
			gen.PrintLN("if err != nil { return }").LN()

			gen.PrintLN("width := 600")
			gen.PrintLN("height := 400")
			gen.PrintLN("if icon.ViewBox.W != 0 || icon.ViewBox.H != 0 {")
			gen.PrintLN("width = int(icon.ViewBox.W)")
			gen.PrintLN("height = int(icon.ViewBox.H)")
			gen.PrintLN("}").LN()

			gen.PrintLN("img = image.NewRGBA(image.Rect(0, 0, width, height))")
			gen.PrintLN("icon.SetTarget(0, 0, float64(width), float64(height))")
			gen.PrintLN("dasher := rasterx.NewDasher(width, height, rasterx.NewScannerGV(width, height, img, img.Bounds()))")
			gen.PrintLN("icon.Draw(dasher, 1.0)")
		},
	)

	// //

	goGen.SeparatorX4().LN()

	goGen.AddStruct("Lang", map[string]generator.GeneratorTypeObj{
		"Info": types("*LangInfoObj", "info"),
		"Sys":  types("*LangSysObj", "sys"),
		"Data": types("*DataObj", "data"),
	})

	// //

	fileName := "treelang_struct.gen.go"
	err = goGen.Save(fileName)
	if err == nil {
		fmt.Printf("The file-struct was created successfully. \n\tDir: %s \n\tFile: %s\n",
			cyan(toDir),
			green(fileName),
		)
	}
	return err
}

//

func writeGoData(fromFilePath, fromReadDir, toDir, packageName string) error {
	fmt.Println("writeGoData", fromFilePath, toDir, packageName)
	return nil
}

func writeGoMap(fromReadDir, toDir, packageName string) error {
	fmt.Println("writeGoMap", toDir, packageName)
	return nil
}

// // // // // //

func init() {
	generateCmd.Flags().StringVar(&CmdFromFilePath, "from", "", "####")
	generateCmd.Flags().StringVar(&CmdToFilePath, "to", "", "####")
	generateCmd.Flags().StringVar(&CmdMasterFile, "master", "", "#####")

	generateCmd.Flags().StringVar(&CmdPackage, "go-package", "", "#####")
	CmdJson = generateCmd.Flags().Bool("json", false, "######")
	CmdYml = generateCmd.Flags().Bool("yml", false, "######")
	CmdMap = generateCmd.Flags().Bool("map", false, "######")

	rootCmd.AddCommand(generateCmd)
}
