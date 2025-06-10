package parser

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
