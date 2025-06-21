package treelang

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// // // // // // // // // // // // // // // // // //

var FileExtension = []string{".yml", ".yaml", ".json"}

type FilePathType byte

const (
	FilePathUnknown FilePathType = iota
	FilePathErr
	FilePathInvalid
	FilePathValid
	FilePathValidDir
	FilePathIsDir
)

func (fp FilePathType) String() string {
	switch fp {
	case FilePathInvalid:
		return "file path is not valid"
	case FilePathValidDir:
		return "file path is valid, but the file does not exist"
	case FilePathErr:
		return "error occurred while accessing the file"
	case FilePathIsDir:
		return "is a folder"
	case FilePathValid:
		return "file exists"
	}
	return "unknown file path state"
}

func CheckFilePath(filePath string) FilePathType {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return FilePathInvalid
	}

	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return FilePathValidDir
	} else if err != nil {
		return FilePathErr
	}

	if info.IsDir() {
		return FilePathIsDir
	}

	return FilePathValid
}

func CheckPathCMD(val *string, param string) (FilePathType, error) {
	if *val == "" {
		return FilePathErr, errors.New("you must specify a file using " + Cyan("--"+param))
	}

	path := CheckFilePath(*val)
	if path != FilePathValid && path != FilePathIsDir && path != FilePathValidDir {
		return path, ParamErr(param, path.String())
	}

	abs, err := filepath.Abs(*val)
	if err != nil {
		return path, ParamErr(param, err.Error())
	}

	if path == FilePathValid {
		ex := filepath.Ext(abs)
		isOK := false

		for _, name := range FileExtension {
			if name == ex {
				isOK = true
				break
			}
		}

		if !isOK {
			return path, ParamErr(param, "Invalid file extension. Must be: "+strings.Join(FileExtension, ", "))
		}
	}

	*val = abs
	return path, nil
}
