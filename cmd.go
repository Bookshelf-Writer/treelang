package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// // // // // // // // // // // // // // // // // //

var rootCmd = &cobra.Command{
	Use:   GlobalName,
	Short: magenta(GlobalName) + " — console utility for working with language trees",
	Long:  magenta(GlobalName) + " — This is a console utility that allows you to create, modify and analyze language trees. You can also generate ready-made files for different programming languages from language trees.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use " + magenta(GlobalName) + " " + cyan("help") + " for more information about a command.")
	},
}

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "Show information about the build",
	Long:    "Show information about the build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s %s\n", magenta(GlobalName), green(GlobalVersion), GlobalHash[:8])
		fmt.Println("Dependencies used:")

		for _, key := range sortMapKey(GlobalDependenciesMap) {
			fmt.Printf("\t%s %s\n", key, blue(GlobalDependenciesMap[key]))
		}
	},
}

// // // // // //

func init() {
	rootCmd.AddCommand(infoCmd)
}
