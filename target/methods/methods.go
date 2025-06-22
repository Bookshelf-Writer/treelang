// Code generated using '_generate/dependencies/main.go'; DO NOT EDIT.
// Generation time: 2025-06-22T01:06:54Z

package methods

import (
	"github.com/Bookshelf-Writer/treelang/methods/diff"
	"github.com/Bookshelf-Writer/treelang/methods/generate"
	"github.com/spf13/cobra"
)

// // // // // // // //

var MethodsMap = map[string]func(*cobra.Command){
	"diff":     diff.Init,
	"generate": generate.Init,
}
