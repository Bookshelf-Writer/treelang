package main

import (
	"encoding/json"
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

	buffer, _ := json.Marshal(obj.Data)
	sys.Hash = Hash(buffer)

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
	obj.Sys = newSys(obj)

	jsonData, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(
		filepath.Join(toDir, "treelang_"+strings.ToLower(obj.Info.Code)+".gen.json"),
		jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createLangGO(obj *LangObj, toDir string) error {
	return nil
}

// //

func createMapYML(arr []*LangInfoObj, toDir string) error {
	maps := make(map[string]*LangInfoObj)
	for _, info := range arr {
		maps[strings.ToLower(info.Code)] = info
	}

	yamlData, err := yaml.Marshal(maps)
	if err != nil {
		return err
	}

	err = os.WriteFile(
		filepath.Join(toDir, "treelang_map.gen.yml"),
		yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createMapJSON(arr []*LangInfoObj, toDir string) error {
	maps := make(map[string]*LangInfoObj)
	for _, info := range arr {
		maps[strings.ToLower(info.Code)] = info
	}

	jsonData, err := json.MarshalIndent(maps, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(
		filepath.Join(toDir, "treelang_map.gen.json"),
		jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createMapGO(arr []*LangInfoObj, toDir string) error {
	return nil
}
