package main

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
)

type GetQuery struct {
	Get string `@"GET"`
	Key string `@Ident`
}

type SetQuery struct {
	Set   string `@"SET"`
	Key   string `@Ident`
	Value string `@String`
}

type DeleteQuery struct {
	Delete string `@"DELETE"`
	Key    string `@Ident`
}

type Query struct {
	Get    *GetQuery    `@@`
	Set    *SetQuery    `| @@`
	Delete *DeleteQuery `| @@`
}

func parseQuery(input string) (*Query, error) {
	parser, err := participle.Build[Query]()
	if err != nil {
		return &Query{}, err
	}

	expr, err := parser.ParseString("", input)
	if err != nil {
		return &Query{}, err
	}
	fmt.Printf("Parsed query: %#v\n", expr)
	return expr, nil
}
