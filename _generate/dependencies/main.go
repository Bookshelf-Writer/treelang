package main

import (
	"bufio"
	_ "embed"
	"fmt"
	dep "github.com/Bookshelf-Writer/treelang"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// // // // // // // // // //

const (
	packageName = "target"
	fileName    = "dependencies.go"
)

//go:embed template.tmpl
var template_text string

//

type Dep struct {
	Key   string
	Value string
}
type TemplateObj struct {
	GenerationTime string
	PackageName    string

	Deps []Dep
}

// //

func main() {
	file, err := os.Open("go.mod")
	if err != nil {
		fmt.Println("An error occurred while opening the file:", err)
		return
	}
	defer file.Close()

	dependencies := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "\t") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				dependencies[fields[0]] = fields[1]
			}
		}
	}

	//

	keys := make([]string, 0, len(dependencies))
	for k := range dependencies {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	depList := make([]Dep, 0, len(keys))
	for _, k := range keys {
		depList = append(depList, Dep{
			Key:   k,
			Value: dependencies[k],
		})
	}

	//

	data := new(TemplateObj)
	data.GenerationTime = time.Now().Format(time.RFC3339)
	data.PackageName = packageName
	data.Deps = depList

	err = dep.WriteFileFromTemplate(filepath.Join("target", fileName), template_text, data)
	if err != nil {
		fmt.Println("An error when trying to save a generated file:", err)
	}
}
