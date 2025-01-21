package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
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
			if masterPath == FilePathUnknown && !*CmdMap {
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
		if masterPath != FilePathValid && !*CmdMap {
			return paramErr("master", errors.New("error with parameter. must point to a file"))
		}

		// //

		if !*CmdMap {
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
	err := writeGoStruct(fromFilePath, toDir, packageName)
	if err == nil {
		return err
	}

	master, slave, err := parseSlave(fromFilePath, fromReadDir)
	if err != nil {
		return err
	}

	for fileName, obj := range slave {
		fmt.Printf("%s:\n", green(fileName))
		finalObj := mergeLangObj(master, obj, 1)
		err = createLangGO(finalObj, toDir)
		if err == nil {
			fmt.Printf("Created: GO %s\n", blue(finalObj.Info.Name.EN))
		}
	}

	return nil
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

// // // // // //

func init() {
	generateCmd.Flags().StringVar(&CmdFromFilePath, "from", "", "path to the folder where the language tree is located. If you specify a file, it will be considered a master-file and the directory where it is located will be a directory with a language tree")
	generateCmd.Flags().StringVar(&CmdToFilePath, "to", "", "path to the folder where to save the generated files")
	generateCmd.Flags().StringVar(&CmdMasterFile, "master", "", "path to the master-file")

	CmdJson = generateCmd.Flags().Bool("json", false, "generate json from language tree")
	CmdYml = generateCmd.Flags().Bool("yml", false, "generate json from language tree")
	CmdMap = generateCmd.Flags().Bool("map", false, "generate a map based on the language tree. Works in conjunction with parameters [json|yml|go-package]")

	generateCmd.Flags().StringVar(&CmdPackage, "go-package", "", "generate go files. with the parameter you need to specify the name of the package")
	CmdGoPNG = generateCmd.Flags().Bool("func-png", false, "Add a predefined method for getting PNG to the generated structure file. Only for [go-package]")

	rootCmd.AddCommand(generateCmd)
}
