.PHONY: run
run:
	go run ./cmd/gourl_shortener

.PHONY: test
test:
	go test ./...

