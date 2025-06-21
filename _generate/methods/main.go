package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	dep "github.com/Bookshelf-Writer/treelang/_generate"
)

// // // // // // // // // //

const (
	packageName = "methods"
	fileName    = "methods.go"

	dirRoot = "methods"
	gitRoot = "github.com/Bookshelf-Writer/treelang"
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
	ImportsArr     []string

	Methods []string
}

// //

func main() {
	methods, err := dirsWithInit(dirRoot)
	if err != nil {
		panic(err)
	}

	fmt.Println(methods)

	for pos, method := range methods {
		method = strings.Split(method, "treelang")[1]
		methods[pos] = gitRoot + method
	}
	sort.Strings(methods)

	panic(errors.New(strings.Join(methods, "\n")))

	//

	data := new(TemplateObj)
	data.GenerationTime = time.Now().Format(time.RFC3339)
	data.PackageName = packageName

	data.Methods = make([]string, 0)
	data.ImportsArr = make([]string, 0)
	data.ImportsArr = append(data.ImportsArr, "github.com/spf13/cobra")

	for _, method := range methods {
		data.ImportsArr = append(data.ImportsArr, method)

		name := strings.Split(method, "/")
		data.Methods = append(data.Methods, name[len(name)-1])
	}

	os.MkdirAll(filepath.Join("target", packageName), os.ModePerm)
	err = dep.WriteFileFromTemplate(filepath.Join("target", packageName, fileName), template_text, data)
	if err != nil {
		fmt.Println("An error when trying to save a generated file:", err)
	}
}
