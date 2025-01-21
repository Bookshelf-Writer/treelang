package main

import (
	"encoding/json"
	"errors"
	"github.com/ghodss/yaml"
	"io"
	"os"
	"path/filepath"
)

// // // // // // // // // // // // // // // // // //

func parseYml(filepath string) (*LangObj, error) {
	yamlFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer yamlFile.Close()

	yamlData, err := io.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}

	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return nil, err
	}

	obj := new(LangObj)
	err = json.Unmarshal(jsonData, obj)
	return obj, err
}

func parseJson(filepath string) (*LangObj, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	obj := new(LangObj)
	err = json.Unmarshal(jsonData, obj)
	return obj, err
}

func ReadFile(adr string) (*LangObj, error) {
	ex := filepath.Ext(adr)
	switch ex {

	case ".yml", ".yaml":
		return parseYml(adr)

	case ".json":
		return parseJson(adr)

	}
	return nil, errors.New("unexpected file type: " + red(ex))
}
