package generate

import (
	"errors"
	"fmt"
	. "github.com/Bookshelf-Writer/treelang"
	"github.com/spf13/cobra"
	"path"
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
				return ParamErr("to", "should only point to a directory")
			}
		}
		if CmdMasterFile != "" {
			masterPath, err = CheckPathCMD(&CmdMasterFile, "master")
			if err != nil {
				return err
			}
			if masterPath != FilePathValid {
				return ParamErr("master", "the path should only point to an existing file")
			}
		}

		if fromPath == FilePathUnknown {
			return ParamErr("from", "this parameter is required")
		}
		if toPath == FilePathUnknown {
			return ParamErr("to", "this parameter is required")
		}

		if fromPath == FilePathIsDir {
			if masterPath == FilePathUnknown && !*CmdMap {
				return ParamErr("master", "this parameter is required")
			}
		}
		if fromPath != FilePathValid && fromPath != FilePathIsDir {
			return ParamErr("from", "error with parameter. must point to a file or folder")
		}
		if fromPath == FilePathValid {
			if CmdMasterFile == "" {
				CmdMasterFile = CmdFromFilePath
				masterPath = FilePathValid
			}
			CmdFromFilePath = path.Dir(CmdFromFilePath)
		}
		if masterPath != FilePathValid && !*CmdMap {
			return ParamErr("master", "error with parameter. must point to a file")
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
			return fmt.Errorf("generation type is not specified. Select %s, %s or enter %s", Cyan("--yml"), Cyan("--json"), Cyan("----go-package"))
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
				return writeGoMap(CmdFromFilePath, CmdToFilePath, CmdMasterFile, CmdPackage)
			} else {
				return writeGoData(CmdMasterFile, CmdFromFilePath, CmdToFilePath, CmdPackage)
			}
		}

		return errors.New("Unexpected termination of context")
	},
}

// // // // // //

func Init(rootCmd *cobra.Command) {
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
