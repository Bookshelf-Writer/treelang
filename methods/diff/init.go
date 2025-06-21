package diff

import (
	"errors"
	. "github.com/Bookshelf-Writer/treelang"
	"github.com/spf13/cobra"
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
				return ParamErr("from", "the path to the file is specified, only the folder is allowed")
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
		if CmdSlaveFile != "" {
			slavePath, err = CheckPathCMD(&CmdSlaveFile, "slave")
			if err != nil {
				return err
			}
			if slavePath != FilePathValid {
				return ParamErr("slave", "the path should only point to an existing file")
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

func Init(rootCmd *cobra.Command) {
	diffCmd.Flags().StringVar(&CmdFromFilePath, "from", "", "the root folder in which the language tree files are located")
	diffCmd.Flags().StringVar(&CmdMasterFile, "master", "", "main file to be compared with")
	diffCmd.Flags().StringVar(&CmdSlaveFile, "slave", "", "specific file to compare with the main file")

	diffCmd.Flags().StringVar(&CmdMode, "mode", "short", "Output mode. By default, outputs in shortened format. [all|short]")
	CmdFull = diffCmd.Flags().Bool("full", false, "full file comparison, including key values")

	rootCmd.AddCommand(diffCmd)
}
