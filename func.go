package main

import (
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/blake2b"
	"sort"
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

func paramErr(param string, err error) error {
	return errors.New(cyan("--"+param) + ":\t" + err.Error())
}

// // // //

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

func Hash(data []byte) string {
	h, _ := blake2b.New(16, []byte(GlobalHash))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
