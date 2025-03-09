timeout=5s
test/tokenizer:
	@go test -timeout ${timeout} -cover ./tokenizer

test/parser:
	@go test -timeout ${timeout} -cover ./parser

test/ast:
	@go test -timeout ${timeout} -cover ./ast

test:
	@make -s test/ast
	@make -s test/tokenizer
	@make -s test/parser

run/demo/def:
	@go run ./demo demo/def.tem

run/demo/def/semicolon:
	@go run ./demo demo/def_semicolon.tem

run/demo/template:
	@go run ./demo demo/template.tem

run/demo/template/semicolon:
	@go run ./demo demo/template_semicolon.tem
