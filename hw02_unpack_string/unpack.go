package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	runeString := []rune(s)
	var resultString strings.Builder
	isMask := false

	for index, char := range runeString {
		if isMask {
			resultString.WriteRune(char)
			isMask = false
			continue
		}

		if char == '\\' {
			isMask = true
			continue
		}

		if num, err := strconv.Atoi(string(char)); err == nil {
			if index == 0 || index < len(runeString)-1 && isDigit(runeString[index+1]) {
				return "", ErrInvalidString
			}
			if num > 0 {
				resultString.WriteString(strings.Repeat(string(runeString[index-1]), num-1))
			} else {
				tmp := []rune(resultString.String())
				resultString.Reset()
				resultString.WriteString(string(tmp[:len(tmp)-1]))
			}
		} else {
			resultString.WriteRune(char)
		}
	}

	return resultString.String(), nil
}

func isDigit(r rune) bool {
	if _, err := strconv.Atoi(string(r)); err == nil {
		return true
	}
	return false
}
