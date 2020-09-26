.PHONY: build clean deploy test

build:
	GOOS=linux GOARCH=amd64 go build -o bin/message ./message
	chmod +x bin/message
	cp message/signup-template.html bin/signup-template.html

clean:
	rm -rf ./bin

deploy: clean build
	npm run deploy

remove:
	npm run remove

test:
	go test -v $$(go list ./... | grep -v /node_modules/)

format:
	gofmt -l -s -w .
