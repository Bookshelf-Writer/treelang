package main

import (
	"bytes"
	"encoding/gob"
	"github.com/ghodss/yaml"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// // // // // // // // // // // // // // // // // //

func newSys(obj *LangObj) *LangSysObj {
	sys := new(LangSysObj)
	sys.Date = time.Now().Format("2006-01-02")

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(obj.Data); err != nil {
		return sys
	}

	sys.Hash = Hash(buffer.Bytes())
	return sys
}

// //

func createLangYML(obj *LangObj, toDir string) error {
	obj.Sys = newSys(obj)

	yamlData, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	err = os.WriteFile(
		filepath.Join(toDir, "treelang_"+strings.ToLower(obj.Info.Code)+".gen.yml"),
		yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createLangJSON(obj *LangObj, toDir string) error {
	return nil
}

func createLangGO(obj *LangObj, toDir string) error {
	return nil
}
