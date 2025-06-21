package treelang

import (
	"encoding/hex"
	"errors"
	"github.com/Bookshelf-Writer/treelang/target"
	"golang.org/x/crypto/blake2b"
	"sort"
	"strings"
	"unicode"
)

// // // // // // // // // // // // // // // // // //

func SortMapKey[T any](mp map[string]T) []string {
	var listBuf []string
	for key := range mp {
		listBuf = append(listBuf, key)
	}
	sort.Strings(listBuf)

	return listBuf
}

func ParamErr(param, err string) error {
	return errors.New(Cyan("--"+param) + ":\t" + err)
}

// // // //

func ToGoVariableName(s string) string {
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
	h, _ := blake2b.New(16, []byte(target.GlobalHash))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
