.PHONY: deps clean build package deloy test

deps:
	go get -t -v ./...

fmt:
	gofmt -s -w bot/*.go

clean:
	rm -rf ./build

test:
	go test -v ./bot/

build: test clean
	GOOS=linux GOARCH=amd64 go build -o build/bot ./bot

run: build
	sam local start-api

package: build
	sam package \
	    --template-file template.yaml \
	    --output-template-file build/packaged.yaml \
	    --s3-bucket filbot

deploy: package
	sam deploy \
	    --template-file build/packaged.yaml \
	    --stack-name filbot \
	    --capabilities CAPABILITY_IAM

status:
	aws cloudformation describe-stacks \
	    --stack-name filbot \
	    --query 'Stacks[].Outputs'
