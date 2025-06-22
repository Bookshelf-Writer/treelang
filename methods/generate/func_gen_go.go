package generate

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/Bookshelf-Writer/treelang"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"
)

// // // // // // // // // // // // // // // // // //

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
	out := make(map[string]map[string][]string)

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

		if _, ok := out[base]; !ok {
			out[base] = make(map[string][]string)
		}

		fullPath := strings.ReplaceAll(fullKey, ":", ".")

		switch val {
		case "str":
			out[base][field] = []string{"string", name, "data." + fullPath}
		case "arr":
			out[base][field] = []string{"[]string", name, "data." + fullPath}
		case "obj":
			out[base][field] = []string{"*" + field + "Obj", name, "data." + fullPath}
		}
	}

	return out
}

// //

const structGoFilename = "treelang_struct.gen"

//go:embed template/template_struct.tmpl
var template_struct string

type TemplateStructObj struct {
	GenerationTime string
	PackageName    string
	ImportsArr     []string

	CmdGoPNG bool
	Struct   map[string]map[string][]string
}

func writeGoStruct(fromFilePath, toDir, packageName string) error {
	data := new(TemplateStructObj)
	data.GenerationTime = time.Now().Format(time.RFC3339)
	data.PackageName = packageName

	data.ImportsArr = make([]string, 0)
	data.CmdGoPNG = *CmdGoPNG

	if *CmdGoPNG {
		data.ImportsArr = append(data.ImportsArr, "bytes")
		data.ImportsArr = append(data.ImportsArr, "image")
		data.ImportsArr = append(data.ImportsArr, "github.com/srwiley/oksvg")
		data.ImportsArr = append(data.ImportsArr, "github.com/srwiley/rasterx")
	}

	// //

	obj, err := ReadFile(fromFilePath)
	if err != nil {
		return err
	}

	data.Struct = structMap(obj.Data)

	//

	fileName := structGoFilename + ".go"
	err = WriteFileFromTemplate(filepath.Join(toDir, fileName), template_struct, data)
	if err == nil {
		fmt.Printf("The file-struct was created successfully. \n\tDir: %s \n\tFile: %s\n",
			Cyan(toDir),
			Green(fileName),
		)
	}
	return err
}

// //

//go:embed template/template_lang.tmpl
var template_lang string

type TemplateLangObj struct {
	GenerationTime string
	PackageName    string

	Name      string
	SmallName string
	DataJson  string
}

func createLangGO(obj *LangObj, toDir, packageName string) error {
	data := new(TemplateLangObj)
	data.GenerationTime = time.Now().Format(time.RFC3339)
	data.PackageName = packageName

	obj.Sys = newSys(obj)
	fileName := genFileName(obj.Info) + ".go"

	//

	data.Name = strings.ToUpper(obj.Info.Code)
	data.SmallName = strings.ToLower(obj.Info.Code)

	jsonData, _ := json.MarshalIndent(obj, "", "  ")
	data.DataJson = string(jsonData)

	//

	err := WriteFileFromTemplate(filepath.Join(toDir, fileName), template_lang, data)
	if err == nil {
		fmt.Printf("The file-struct was created successfully. \n\tDir: %s \n\tFile: %s\n",
			Cyan(toDir),
			Green(fileName),
		)
	}
	return err
}

// //

//go:embed template/template_map.tmpl
var template_map string

type TemplateMapObj struct {
	GenerationTime string
	PackageName    string

	Langs      []string
	MasterLang string
}

func createMapGO(arr []*LangInfoObj, master *LangInfoObj, toDir, packageName string) error {
	if len(arr) == 0 {
		return errors.New("can't build a map, no data")
	}

	data := new(TemplateMapObj)
	data.GenerationTime = time.Now().Format(time.RFC3339)
	data.PackageName = packageName

	//

	data.Langs = make([]string, 0)
	for _, info := range arr {
		data.Langs = append(data.Langs, strings.ToUpper(info.Code))
	}
	sort.Strings(data.Langs)

	data.MasterLang = strings.ToUpper(master.Code)
	for i, val := range data.Langs {
		if val == data.MasterLang {
			if i == 0 {
				break
			}

			data.Langs[0], data.Langs[i] = data.Langs[i], data.Langs[0]
			break
		}
	}

	//

	err := WriteFileFromTemplate(filepath.Join(toDir, genMapName+".go"), template_map, data)
	if err == nil {
		fmt.Printf("The file-map was created successfully. \n\tDir: %s \n\tFile: %s\n",
			Cyan(toDir),
			Green(genMapName+".go"),
		)
	}
	return err
}
