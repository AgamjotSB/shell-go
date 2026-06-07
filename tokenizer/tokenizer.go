// Package tokenizer provides a tokenizer to split lines into tokens
package tokenizer

func Parse(input string) ([]string, error) {
	lexer := NewLexer(input)

	var tokens []string

	for {
		token, isEOF, err := lexer.NextToken()
		if err != nil {
			return nil, err
		}

		if isEOF {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}
