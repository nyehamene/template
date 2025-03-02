test/tokenizer:
	@go test -timeout 5s -cover ./tokenizer

test:
	@make test/tokenizer

run/demo/def:
	@go run ./demo demo/def.tem

run/demo/def/semicolon:
	@go run ./demo demo/def_semicolon.tem

run/demo/template:
	@go run ./demo demo/template.tem

run/demo/template/semicolon:
	@go run ./demo demo/template_semicolon.tem
