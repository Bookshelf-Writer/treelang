package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

// // // // // // // // // // // // // // // // // //

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "сравнение файлов",
	RunE: func(cmd *cobra.Command, args []string) error {
		var fromPath, toPath, masterPath, slavePath FilePathType
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
		if CmdToFilePath != "" {
			toPath, err = CheckPathCMD(&CmdToFilePath, "to")
			if err != nil {
				return err
			}
			if toPath == FilePathValid || toPath == FilePathValidDir {
				return paramErr("to", errors.New("the path to the file is specified, only the folder is allowed"))
			}
		}
		if CmdMasterFile != "" {
			masterPath, err = CheckPathCMD(&CmdMasterFile, "master")
			if err != nil {
				return err
			}
		}
		if CmdSlaveFile != "" {
			slavePath, err = CheckPathCMD(&CmdSlaveFile, "slave")
			if err != nil {
				return err
			}
		}

		fmt.Println(fromPath, toPath, masterPath, slavePath)
		return nil
	},
}

// // // // // //

func init() {
	diffCmd.Flags().StringVar(&CmdMasterFile, "master", "", "######")
	diffCmd.Flags().StringVar(&CmdSlaveFile, "slave", "", "######")

	diffCmd.Flags().StringVar(&CmdFromFilePath, "from", "", "######")
	diffCmd.Flags().StringVar(&CmdToFilePath, "to", "", "######")

	rootCmd.AddCommand(diffCmd)
}
