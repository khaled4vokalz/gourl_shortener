.PHONY: run-server test-server run-client

run-server:
	cd server; \
	go run ./cmd/gourl_shortener

test-server:
	cd server; \
	go test ./...

run-client:
	cd client; \
	npm start

