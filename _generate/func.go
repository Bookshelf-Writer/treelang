package _generate

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// // // // // // // // // //

func WriteFileFromTemplate(pathToFile string, textTemplate string, dataTemplate any) error {
	fileName := filepath.Base(pathToFile)

	tmpl := template.New("cli-template").Funcs(template.FuncMap{
		"split": strings.Split,
	})

	t, err := tmpl.New(fileName).Parse(textTemplate)
	if err != nil {
		return fmt.Errorf("init template [%s]: %s", fileName, err.Error())
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, dataTemplate)
	if err != nil {
		return fmt.Errorf("filling template [%s]: %s", fileName, err.Error())
	}

	file, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open file [%s]: %s", fileName, err.Error())
	}
	defer file.Close()

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format template [%s]: %s", fileName, err.Error())
	}

	_, err = file.Write(formatted)
	if err != nil {
		return fmt.Errorf("write file [%s]: %s", fileName, err.Error())
	}

	fmt.Println("\tGenerate: " + pathToFile)
	return nil
}
