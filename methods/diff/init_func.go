package diff

import (
	"fmt"
	. "github.com/Bookshelf-Writer/treelang"
	"os"
	"path/filepath"
	"strings"
)

// // // // // // // // // // // // // // // // // //

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
			fmt.Printf("%s ", Yellow(name))

			if obj.Info.Name == nil {
				fmt.Printf("%s\t", Red("NAME-UNDEF"))

			} else {
				if obj.Info.Code == "" {
					fmt.Printf("%s ", Red("CODE-UNDEF"))
				} else {
					fmt.Printf("%s ", Blue(strings.ToUpper(obj.Info.Code)))
				}
				fmt.Printf("%s\t", Green(obj.Info.Name.DEF))
			}
		}
		return nil
	} else {
		for name, obj := range load {
			fmt.Printf("%s\n", Yellow(name))

			if obj.Info.Name == nil {
				fmt.Printf("\tNAME: %s\n", Red("\tUNDEF"))

			} else {
				if obj.Info.Code == "" {
					fmt.Printf("\tCODE: %s\n", Red("\tUNDEF"))
				} else {
					fmt.Printf("\t%s [%s]\n", Blue(strings.ToUpper(obj.Info.Code)), Magenta(obj.Info.Name.EN))
				}

				if obj.Info.Name.DEF == "" {
					fmt.Printf("\tDefName: %s\n", Red("\tUNDEF"))
				} else {
					fmt.Printf("\t%s\n", Cyan(obj.Info.Name.DEF))
				}
			}

			if obj.Info.Flag == "" {
				fmt.Printf("\tFLAG: %s\n", Red("\tUNDEF"))
			} else {
				fmt.Printf("\tFLAG: %s\n", Green("\tOK"))
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
		arr, err = DiffFile(master, slave)
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

		arr, err = DiffObj(masterFile, slaveFile)
		if err != nil {
			return err
		}
	}

	if len(arr) == 0 {
		fmt.Printf("No differences between %s and %s\n", Green(filepath.Base(master)), Magenta(filepath.Base(slave)))
	} else {
		fmt.Printf("%s >> %s\n", Green(filepath.Base(master)), Magenta(filepath.Base(slave)))
	}

	for _, txt := range arr {
		for _, line := range strings.Split(txt, "\n") {
			if line == "" {
				continue
			}

			switch strings.SplitAfter(line, "")[0] {
			case "+":
				fmt.Println(strings.Repeat("|\t", pad) + Red(line))
			case "-":
				fmt.Println(strings.Repeat("|\t", pad) + Yellow(line))
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
			for _, name := range FileExtension {
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
