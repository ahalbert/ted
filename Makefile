
build:
	go build -o bin/fsaed cmd/main.go
install:
	go build -o bin/fsaed cmd/main.go
	cp bin/fsaed $$(go env GOPATH)/bin
clean:
	rm -rf bin
	rm -rf program.fsa
	rm -rf test
test: build
	./tests/test.zsh
