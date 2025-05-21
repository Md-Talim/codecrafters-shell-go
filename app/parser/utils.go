package parser

import (
	"strings"

	"github.com/md-talim/codecrafters-shell-go/app/shellio"
)

func mapBackshlash(character byte) byte {
	if character == DOUBLE || character == BACKSLASH {
		return character
	}
	return END
}

func isRedirectionOperator(operator string) bool {
	return (operator == ">") || (operator == "1>") || (operator == "2>") ||
		(operator == ">>") || (operator == "1>>") || (operator == "2>>")
}

func getRedirectionConfig(operator string, file string) shellio.RedirectionConfig {
	descriptor := 1
	append := false
	if strings.HasPrefix(operator, "2") {
		descriptor = 2
	}
	if strings.Contains(operator, ">>") {
		append = true
	}
	return shellio.RedirectionConfig{
		File:            file,
		Descriptor:      descriptor,
		IsAppendEnabled: append,
	}
}
