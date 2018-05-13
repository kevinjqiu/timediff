.PHONY: test

test:
	go test -v

coverage:
	go test -coverprofile=cover.out
	go tool cover -func cover.out
