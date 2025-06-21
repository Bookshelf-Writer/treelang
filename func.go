package main

import (
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/blake2b"
	"sort"
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

func paramErr(param string, err error) error {
	return errors.New(cyan("--"+param) + ":\t" + err.Error())
}

// // // //

// ToGoGlobalName перетворює довільний рядок на коректний експортований
// ідентифікатор у стилі PascalCase (CamelCase з великої літери).
// Правила:
// /  • усі пробіли та будь-які неалфавітно-цифрові символи вважаються роздільниками;
// /  • усе, що йде після роздільника, починається з великої літери;
// /  • якщо результат починається не з літери (наприклад, з цифри), додаємо префікс "X".
func ToGoGlobalName(s string) string {
	// Розбиваємо рядок за будь-яким символом, який НЕ є літерою чи цифрою.
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	var b strings.Builder
	for _, w := range words {
		if w == "" {
			continue
		}
		rs := []rune(w)
		b.WriteRune(unicode.ToUpper(rs[0]))
		for _, r := range rs[1:] {
			b.WriteRune(unicode.ToLower(r))
		}
	}

	name := b.String()
	if name == "" || !unicode.IsLetter([]rune(name)[0]) {
		name = "X" + name
	}

	return name
}

func Hash(data []byte) string {
	h, _ := blake2b.New(16, []byte(GlobalHash))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
