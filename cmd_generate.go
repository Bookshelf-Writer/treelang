package main

import (
	"errors"
	"fmt"
	generator "github.com/Bookshelf-Writer/SimpleGenerator"
	"github.com/spf13/cobra"
	"path/filepath"
	"reflect"
	"strings"
)

// // // // // // // // // // // // // // // // // //

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generating files from a tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		var fromPath, masterPath, slavePath, toPath FilePathType
		var err error

		if CmdFromFilePath != "" {
			fromPath, err = CheckPathCMD(&CmdFromFilePath, "from")
			if err != nil {
				return err
			}
		}
		if CmdToFilePath != "" {
			toPath, err = CheckPathCMD(&CmdToFilePath, "from")
			if err != nil {
				return err
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
		if CmdSlaveFile != "" {
			slavePath, err = CheckPathCMD(&CmdSlaveFile, "slave")
			if err != nil {
				return err
			}
			if slavePath != FilePathValid {
				return paramErr("slave", errors.New("the path should only point to an existing file"))
			}
		}

		fmt.Println(fromPath, masterPath, slavePath, toPath)

		writeStructGo(filepath.Join(CmdFromFilePath, "en.json"), "main")

		return nil
	},
}

// // // // // //

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

func writeStructGo(fromFilePath, packageName string) error {
	goGen := generator.NewFilePathName(filepath.Dir(fromFilePath), packageName)
	types := func(name, jsonName string) generator.GeneratorTypeObj {
		return generator.GeneratorTypeObj{Types: goGen.NewType(name), Tags: map[string]string{"json": jsonName, "xml": jsonName, "yaml": jsonName}}
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

	err = goGen.Save("treelang_struct.gen.go")
	if err == nil {
		fmt.Printf("The file-struct was created successfully. \n\tDir: %s \n\tFile: %s\n",
			cyan(filepath.Dir(fromFilePath)),
			green("treelang_struct.gen.go"),
		)
	}
	return err
}

// // // // // //

func init() {
	generateCmd.Flags().StringVar(&CmdFromFilePath, "from", "", "####")
	generateCmd.Flags().StringVar(&CmdToFilePath, "to", "", "####")

	generateCmd.Flags().StringVar(&CmdMasterFile, "master", "", "#####")
	generateCmd.Flags().StringVar(&CmdSlaveFile, "slave", "", "#####")

	generateCmd.Flags().StringVar(&CmdPackage, "go-package", "", "#####")
	CmdJson = diffCmd.Flags().Bool("json", false, "######")
	CmdYml = diffCmd.Flags().Bool("yml", false, "######")

	CmdMap = diffCmd.Flags().Bool("map", false, "######")
	CmdLangs = diffCmd.Flags().Bool("langs", false, "######")

	rootCmd.AddCommand(generateCmd)
}
