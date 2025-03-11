timeout=5s

test/tokenizer:
	@go test -timeout ${timeout} -cover ./tokenizer

test/parser:
	@go test -timeout ${timeout} -cover ./parser

test/ast:
	@go test -timeout ${timeout} -cover ./ast

test/demo:
	@go test -timeout ${timeout} ./demo

test/queue:
	@go test -timeout ${timeout} -cover ./queue

test:
	@make -s test/ast
	@make -s test/tokenizer
	@make -s test/parser
	@make -s test/queue

run/demo:
	@go run ./demo
