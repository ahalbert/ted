
build:
	go build -o bin/ted ted.go
install:
	go install ted.go
clean:
	rm -rf bin
	rm -rf program.fsa
	rm -rf test
test: build
	./tests/test.zsh
