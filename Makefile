timeout=5s
test/tokenizer:
	@go test -timeout ${timeout} -cover ./tokenizer

test/parser:
	@go test -timeout ${timeout} -cover ./parser

test/queue:
	@go test -timeout ${timeout} -cover ./queue

test:
	@make test/tokenizer
	@make test/parser

run/demo/def:
	@go run ./demo demo/def.tem

run/demo/def/semicolon:
	@go run ./demo demo/def_semicolon.tem

run/demo/template:
	@go run ./demo demo/template.tem

run/demo/template/semicolon:
	@go run ./demo demo/template_semicolon.tem
