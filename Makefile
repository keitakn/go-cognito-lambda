.PHONY: build clean deploy test

build:
	GOOS=linux GOARCH=amd64 go build -o bin/message ./message
	chmod +x bin/message

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

test:
	go test -v ./...
