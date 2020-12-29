.PHONY: build clean deploy test

build:
	GOOS=linux GOARCH=amd64 go build -o bin/message ./message
	chmod +x bin/message
	cp message/signup-template.html bin/signup-template.html
	cp message/forgot-password-template.html bin/forgot-password-template.html
	GOOS=linux GOARCH=amd64 go build -o bin/changepassword ./api/changepassword/main.go
	chmod +x bin/changepassword
	GOOS=linux GOARCH=amd64 go build -o bin/defineauthchallenge ./authchallenge/define/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/createauthchallenge ./authchallenge/create/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/verifyauthchallenge ./authchallenge/verify/main.go

clean:
	rm -rf ./bin

deploy: clean build
	npm run deploy

remove:
	npm run remove

test:
	go clean -testcache
	go test -v $$(go list ./... | grep -v /node_modules/)

format:
	gofmt -l -s -w .
	goimports -w -l ./
