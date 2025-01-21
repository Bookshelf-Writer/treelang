package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

// // // // // // // // // // // // // // // // // //

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "displaying information about files and comparing them",
	RunE: func(cmd *cobra.Command, args []string) error {
		var fromPath, masterPath, slavePath FilePathType
		var err error

		if CmdFromFilePath != "" {
			fromPath, err = CheckPathCMD(&CmdFromFilePath, "from")
			if err != nil {
				return err
			}
			if fromPath == FilePathValid || fromPath == FilePathValidDir {
				return paramErr("from", errors.New("the path to the file is specified, only the folder is allowed"))
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

		if masterPath == FilePathUnknown && slavePath == FilePathUnknown && fromPath == FilePathIsDir {
			return checkFolder(CmdFromFilePath)
		}

		if masterPath == FilePathValid && slavePath == FilePathUnknown && fromPath == FilePathIsDir {
			return diffMasterToDir(CmdMasterFile, CmdFromFilePath)
		}

		if masterPath == FilePathValid && slavePath == FilePathValid {
			return diffMasterToSlave(CmdMasterFile, CmdSlaveFile, 0)
		}

		return errors.New("you should use parameters")
	},
}

// // // // // //

func checkFolder(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	load := make(map[string]*LangObj)
	for _, file := range files {
		if !file.IsDir() {
			obj, err := ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				continue
			}

			if obj.Info != nil && obj.Data != nil {
				load[file.Name()] = obj
			}
		}
	}

	if CmdMode == "short" {
		for name, obj := range load {
			fmt.Printf("%s ", yellow(name))

			if obj.Info.Name == nil {
				fmt.Printf("%s\t", red("NAME-UNDEF"))

			} else {
				if obj.Info.Code == "" {
					fmt.Printf("%s ", red("CODE-UNDEF"))
				} else {
					fmt.Printf("%s ", blue(strings.ToUpper(obj.Info.Code)))
				}
				fmt.Printf("%s\t", green(obj.Info.Name.DEF))
			}
		}
		return nil
	} else {
		for name, obj := range load {
			fmt.Printf("%s\n", yellow(name))

			if obj.Info.Name == nil {
				fmt.Printf("\tNAME: %s\n", red("\tUNDEF"))

			} else {
				if obj.Info.Code == "" {
					fmt.Printf("\tCODE: %s\n", red("\tUNDEF"))
				} else {
					fmt.Printf("\t%s [%s]\n", blue(strings.ToUpper(obj.Info.Code)), magenta(obj.Info.Name.EN))
				}

				if obj.Info.Name.DEF == "" {
					fmt.Printf("\tDefName: %s\n", red("\tUNDEF"))
				} else {
					fmt.Printf("\t%s\n", cyan(obj.Info.Name.DEF))
				}
			}

			if obj.Info.Flag == "" {
				fmt.Printf("\tFLAG: %s\n", red("\tUNDEF"))
			} else {
				fmt.Printf("\tFLAG: %s\n", green("\tOK"))
			}

			if obj.Sys != nil {
				if obj.Sys.Date != "" && obj.Sys.Hash != "" {
					fmt.Printf("\t%s\t%s\n", obj.Sys.Date, obj.Sys.Hash[:8])
				}
			}

		}

		return nil
	}
}

func diffMasterToSlave(master, slave string, pad int) error {
	var arr []string
	var err error

	if *CmdFull {
		arr, err = diffFile(master, slave)
		if err != nil {
			return err
		}
	} else {
		masterFile, err := ReadFile(master)
		if err != nil {
			return err
		}
		slaveFile, err := ReadFile(slave)
		if err != nil {
			return err
		}

		arr, err = diffObj(masterFile, slaveFile)
		if err != nil {
			return err
		}
	}

	if len(arr) == 0 {
		fmt.Printf("No differences between %s and %s\n", green(filepath.Base(master)), magenta(filepath.Base(slave)))
	} else {
		fmt.Printf("%s >> %s\n", green(filepath.Base(master)), magenta(filepath.Base(slave)))
	}

	for _, txt := range arr {
		for _, line := range strings.Split(txt, "\n") {
			if line == "" {
				continue
			}

			switch strings.SplitAfter(line, "")[0] {
			case "+":
				fmt.Println(strings.Repeat("|\t", pad) + red(line))
			case "-":
				fmt.Println(strings.Repeat("|\t", pad) + yellow(line))
			default:
				fmt.Println(strings.Repeat("|\t", pad) + line)
			}

		}
		fmt.Println(strings.Repeat("|\t", pad) + "")
	}
	return nil
}

func diffMasterToDir(master, path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			ex := filepath.Ext(file.Name())
			for _, name := range fileExtension {
				if name == ex {
					if filepath.Base(master) != file.Name() {
						diffMasterToSlave(master, filepath.Join(path, file.Name()), 1)
						fmt.Println("")
					}
					break
				}
			}
		}
	}
	return nil
}

// // // // // //

func init() {
	diffCmd.Flags().StringVar(&CmdFromFilePath, "from", "", "the root folder in which the language tree files are located")
	diffCmd.Flags().StringVar(&CmdMasterFile, "master", "", "main file to be compared with")
	diffCmd.Flags().StringVar(&CmdSlaveFile, "slave", "", "specific file to compare with the main file")

	diffCmd.Flags().StringVar(&CmdMode, "mode", "short", "Output mode. By default, outputs in shortened format. [all|short]")
	CmdFull = diffCmd.Flags().Bool("full", false, "full file comparison, including key values")

	rootCmd.AddCommand(diffCmd)
}
