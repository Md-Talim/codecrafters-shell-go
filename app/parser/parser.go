package parser

import "strings"

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

func (p *Parser) Parse() []string {
	var arguments []string

	for {
		argument := p.nextArgument()
		if argument == nil {
			break
		}

		arguments = append(arguments, *argument)
	}

	return arguments
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
			p.handleBackshalsh(&builder)
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
				builder.WriteByte(character)
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

func (p *Parser) handleBackshalsh(builder *strings.Builder) {
	character := p.next()
	if character == END {
		return
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
