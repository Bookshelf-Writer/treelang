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
