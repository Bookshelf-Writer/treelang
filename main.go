//go:generate bash ./_run/creator_const_Go.sh
//go:generate bash ./_run/creator_dependencies_Go.sh

package main

import (
	"fmt"
	"os"
)

// // // // // // // // // // // // // // // // // //

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
