.PHONY: unit-tests tidy
unit-tests: tidy
	CGO_ENABLED=0 go test -count=1 -v .

tidy:
	go mod tidy
