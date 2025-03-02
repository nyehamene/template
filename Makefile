test/tokenizer:
	@go test -timeout 5s -cover ./tokenizer

test:
	@make test/tokenizer

demo/def:
	@go run ./demo demo/def.tem

