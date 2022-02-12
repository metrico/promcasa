package labels

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

func ParseLabels(labels string) (map[string]string, error) {
	s := scanner.Scanner{}
	s.Init(strings.NewReader(labels))
	errorF := func() (map[string]string, error) {
		return nil, fmt.Errorf("Unknown input: %s", labels[s.Offset:])
	}
	tok := s.Scan()
	checkRune := func(expect rune, strExpect string) bool {
		return tok == expect && (strExpect == "" || s.TokenText() == strExpect)
	}
	if !checkRune(123, "{") {
		return errorF()
	}
	res := map[string]string{}
	for tok != scanner.EOF {
		tok = s.Scan()
		if !checkRune(scanner.Ident, "") {
			return errorF()
		}
		name := s.TokenText()
		tok = s.Scan()
		if !checkRune(61, "=") {
			return errorF()
		}
		tok = s.Scan()
		if !checkRune(scanner.String, "") {
			return errorF()
		}
		val, err := strconv.Unquote(s.TokenText())
		if err != nil {
			return nil, err
		}
		tok = s.Scan()
		res[name] = val
		if checkRune(125, "}") {
			return res, nil
		}
		if !checkRune(44, ",") {
			return errorF()
		}
	}
	return res, nil
}
