package main

import (
	"fmt"
	. "github.com/Bookshelf-Writer/treelang"
	"github.com/Bookshelf-Writer/treelang/target"
	"github.com/Bookshelf-Writer/treelang/target/methods"
	"github.com/spf13/cobra"
)

// // // // // // // // // // // // // // // // // //

var rootCmd = &cobra.Command{
	Use:   target.GlobalName,
	Short: Magenta(target.GlobalName) + " — console utility for working with language trees",
	Long:  Magenta(target.GlobalName) + " — This is a console utility that allows you to create, modify and analyze language trees. You can also generate ready-made files for different programming languages from language trees.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use " + Magenta(target.GlobalName) + " " + Cyan("help") + " for more information about a command.")
	},
}

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i"},
	Short:   "Show information about the build",
	Long:    "Show information about the build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s %s\n", Magenta(target.GlobalName), Green(target.GlobalVersion), target.GlobalHash[:8])
		fmt.Println("Dependencies used:")

		for _, key := range SortMapKey(target.GlobalDependenciesMap) {
			fmt.Printf("\t%s %s\n", key, Blue(target.GlobalDependenciesMap[key]))
		}
	},
}

// // // // // //

func init() {
	rootCmd.AddCommand(infoCmd)

	for _, f := range methods.MethodsMap {
		f(rootCmd)
	}
}
