.PHONY: build helloworld

build:
	sam build

helloworld
	GOOS=linux GOARCH=amd64 go build -o bin/helloworld ./helloworld
	chmod +x bin/helloworld
