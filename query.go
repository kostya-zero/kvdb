package main

import (
	"github.com/alecthomas/participle/v2"
)

type Location struct {
	Db  string `@Ident`
	Dot string `@"."`
	Key string `@Ident`
}

type CreateDbQuery struct {
	CreateDb string `@"CREATEDB"`
	Name     string `@Ident`
}

type GetQuery struct {
	Get      string    `@"GET"`
	Location *Location `@@`
}

type SetQuery struct {
	Set      string    `@"SET"`
	Location *Location `@@`
	Value    string    `@String`
}

type RemoveQuery struct {
	Delete string  `@"REMOVE"`
	Which  string  `@("DB" | "KEY")`
	DB     string  `@Ident`
	Key    *string `( "." @Ident )?`
}

type UpdateQuery struct {
	Update   string    `@"UPDATE"`
	Location *Location `@@`
	Value    string    `@String`
}

type Query struct {
	CreateDb *CreateDbQuery `@@`
	Get      *GetQuery      `| @@`
	Set      *SetQuery      `| @@`
	Remove   *RemoveQuery   `| @@`
	Update   *UpdateQuery   `| @@`
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
	return expr, nil
}
