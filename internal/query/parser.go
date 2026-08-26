package query

import (
	"fmt"
	"github.com/example/distributed-property-graph/internal/platform"
	"strconv"
	"strings"
)

type Token struct {
	Kind  string
	Value string
}

func Tokenize(input string) ([]Token, error) {
	fields := strings.Fields(input)
	out := []Token{}
	for _, field := range fields {
		switch {
		case field == "->":
			out = append(out, Token{"arrow", field})
		case strings.HasPrefix(field, "depth="):
			if _, err := strconv.Atoi(strings.TrimPrefix(field, "depth=")); err != nil {
				return nil, platform.ErrInvalid
			}
			out = append(out, Token{"depth", strings.TrimPrefix(field, "depth=")})
		case strings.HasPrefix(field, "type="):
			out = append(out, Token{"type", strings.TrimPrefix(field, "type=")})
		default:
			return nil, fmt.Errorf("unknown token %s: %w", field, platform.ErrInvalid)
		}
	}
	return out, nil
}
func Parse(input string) (Request, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return Request{}, err
	}
	req := Request{}
	for _, token := range tokens {
		switch token.Kind {
		case "depth":
			req.Depth, _ = strconv.Atoi(token.Value)
		case "type":
			req.EdgeType = token.Value
		}
	}
	return req, nil
}
