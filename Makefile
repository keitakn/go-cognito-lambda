.PHONY: build clean deploy test

build:
	GOOS=linux GOARCH=amd64 go build -o bin/message ./message
	chmod +x bin/message

clean:
	rm -rf ./bin

deploy: clean build
	npm run deploy

remove:
	npm run remove

test:
	go test -v ./...

format:
	gofmt -l -s -w .
