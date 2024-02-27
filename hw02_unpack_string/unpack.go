package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidString is an error indicating an invalid string.
var ErrInvalidString = errors.New("invalid string")

// Unpack parses a string and performs unpacking.
func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	runes := make([]rune, 0, len(input))
	// fmt.Printf("runes %v, %v, %v\n", runes, cap(runes), len(runes))
	isLastCharacterDigit := false
	var lastCharacter rune
	for _, r := range input {
		digit, e := strconv.Atoi(string(r))
		if e != nil { // letter
			lastCharacter = r
			runes = append(runes, r)
			isLastCharacterDigit = false
			continue
		}

		// digit
		if lastCharacter == 0 { // first element is digit
			return "", ErrInvalidString
		}
		if isLastCharacterDigit { // there are a few digits near
			return "", ErrInvalidString
		}
		isLastCharacterDigit = true

		if digit == 0 {
			if len(runes) > 0 {
				runes = runes[:len(runes)-1]
			}
		} else if digit > 1 {
			tail := strings.Repeat(string(lastCharacter), digit-1)
			runes = append(runes, []rune(tail)...)
		}
	}
	// fmt.Printf("runes %q, %v, %v\n\n", runes, cap(runes), len(runes))
	result := string(runes)
	return result, nil
}
