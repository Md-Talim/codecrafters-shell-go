package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/md-talim/codecrafters-shell-go/internal/shellio"
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

func (p *Parser) Parse() ([][]string, shellio.RedirectionConfig) {
	var (
		allCommands        [][]string
		currentCommandArgs []string
		redirection        shellio.RedirectionConfig
	)

	for {
		argument := p.nextArgument()
		if argument == nil {
			if len(currentCommandArgs) > 0 {
				allCommands = append(allCommands, currentCommandArgs)
			}
			break
		}

		token := *argument

		if token == "|" {
			if len(currentCommandArgs) == 0 {
				return nil, shellio.RedirectionConfig{}
			}
			allCommands = append(allCommands, currentCommandArgs)
			currentCommandArgs = []string{} // Reset for the next command
		} else if isRedirectionOperator(*argument) {
			// Redirection applies to the ouput of the last command in the pipeline.
			// Expect this at the end of the command line.
			if len(allCommands) == 0 && len(currentCommandArgs) == 0 {
				return nil, shellio.RedirectionConfig{}
			}
			// Add pending arguments for the current (which will be the last) command.
			if len(currentCommandArgs) > 0 {
				allCommands = append(allCommands, currentCommandArgs)
			}
			// If currentCommandArgs is empt here, but allCommands is not, it implies "cmd1 | > out".
			// This is a syntax error (missing command after pipe before redirection).
			if len(allCommands) > 0 && len(allCommands[len(allCommands)-1]) == 0 {
				return nil, shellio.RedirectionConfig{}
			}

			fileName := p.nextArgument()
			if fileName == nil {
				fmt.Fprintln(os.Stderr, "Error: Missing file name for redirection")
				return nil, shellio.RedirectionConfig{}
			}
			redirection = getRedirectionConfig(token, *fileName)
			break
		} else {
			currentCommandArgs = append(currentCommandArgs, token)
		}
	}
	return allCommands, redirection
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
