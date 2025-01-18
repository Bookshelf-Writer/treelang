package main

import (
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/blake2b"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// // // // // // // // // // // // // // // // // //

func sortMapKey[T any](mp map[string]T) []string {
	var listBuf []string
	for key := range mp {
		listBuf = append(listBuf, key)
	}
	sort.Strings(listBuf)

	return listBuf
}

func sortMapKeyInt[T any](mp map[string]T) []string {
	var listBuf []int
	for key := range mp {
		value, _ := strconv.Atoi(key)
		listBuf = append(listBuf, value)
	}
	sort.Ints(listBuf)

	var listBufString []string
	for _, value := range listBuf {
		listBufString = append(listBufString, strconv.Itoa(value))
	}
	return listBufString
}

func stringAdr(arr []bool) string {
	var str string
	for _, b := range arr {
		if b {
			str += "1"
		} else {
			str += "0"
		}
	}
	return str
}

func paramErr(param string, err error) error {
	return errors.New(cyan("--"+param) + ":\t" + err.Error())
}

// //

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
		return FilePathErr, errors.New("you must specify a file using " + cyan("--"+param))
	}

	path := CheckFilePath(*val)
	if path != FilePathValid && path != FilePathIsDir && path != FilePathValidDir {
		return path, paramErr(param, errors.New(path.String()))
	}

	abs, err := filepath.Abs(*val)
	if err != nil {
		return path, paramErr(param, err)
	}

	if path == FilePathValid {
		ex := filepath.Ext(abs)
		isOK := false

		for _, name := range fileExtension {
			if name == ex {
				isOK = true
				break
			}
		}

		if !isOK {
			return path, paramErr(param, errors.New("Invalid file extension. Must be: "+strings.Join(fileExtension, ", ")))
		}
	}

	*val = abs
	return path, nil
}

// //

func ToGoVariableName(input string) string {
	var result []rune
	capitalizeNext := true

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if capitalizeNext {
				result = append(result, unicode.ToUpper(r))
				capitalizeNext = false
			} else {
				result = append(result, r)
			}
		} else if unicode.IsSpace(r) {
			capitalizeNext = true
		}
	}

	return string(result)
}

func Hash(content string) string {
	h, _ := blake2b.New(16, []byte(GlobalHash))
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}
