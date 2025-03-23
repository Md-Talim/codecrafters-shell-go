package parser

import (
	"log"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/app/shellio"
)

const (
	END       = '\x00' // Null character
	SPACE     = ' '    // Space character
	SINGLE    = '\''   // Single quote
	DOUBLE    = '"'    // Double quote
	BACKSLASH = '\\'   // Backslash
)

type Parser struct {
	Input string
	Index int
}

func NewParser(input string) Parser {
	return Parser{
		Input: input,
		Index: -1,
	}
}

func (p *Parser) Parse() ([]string, shellio.RedirectionConfig) {
	var arguments []string
	var redirection shellio.RedirectionConfig

	for {
		argument := p.nextArgument()
		if argument == nil {
			break
		}
		if isRedirectionOperator(*argument) {
			file := p.nextArgument()
			if file == nil {
				log.Fatalln("Expect file name after >")
				break
			}

			if strings.HasPrefix(*argument, "2") {
				redirection = shellio.RedirectionConfig{File: *file, Descriptor: 2}
			} else {
				redirection = shellio.RedirectionConfig{File: *file, Descriptor: 1}
			}
			break
		}
		arguments = append(arguments, *argument)
	}

	return arguments, redirection
}

func (p *Parser) nextArgument() *string {
	builder := strings.Builder{}

	for {
		character := p.next()
		if character == END {
			break
		}

		switch character {
		case SPACE:
			if builder.Len() > 0 {
				result := builder.String()
				return &result
			}
		case BACKSLASH:
			p.handleBackshalsh(&builder, false)
		case SINGLE:
			for {
				character = p.next()
				if character == END || character == SINGLE {
					break
				}
				builder.WriteByte(character)
			}
		case DOUBLE:
			for {
				character = p.next()
				if character == END || character == DOUBLE {
					break
				}
				if character == BACKSLASH {
					p.handleBackshalsh(&builder, true)
				} else {
					builder.WriteByte(character)
				}
			}
		default:
			builder.WriteByte(character)
		}
	}

	if builder.Len() > 0 {
		result := builder.String()
		return &result
	}

	return nil
}

func (p *Parser) handleBackshalsh(builder *strings.Builder, inQuotes bool) {
	character := p.next()
	if character == END {
		return
	}
	if inQuotes {
		mapped := mapBackshlash(character)
		if mapped != END {
			character = mapped
		} else {
			builder.WriteByte(BACKSLASH)
		}
	}
	builder.WriteByte(character)
}

func (p *Parser) next() byte {
	p.Index++
	if p.Index >= len(p.Input) {
		return END
	}

	return p.Input[p.Index]
}
