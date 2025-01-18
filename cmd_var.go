package main

// // // // // // // // // // // // // // // // // //

var (
	fileExtension = []string{".yml", ".yaml", ".json"}

	CmdMasterFile   string
	CmdSlaveFile    string
	CmdFromFilePath string
	CmdToFilePath   string

	CmdMode    string // [all|short]
	CmdFull    *bool
	CmdPackage string
)
