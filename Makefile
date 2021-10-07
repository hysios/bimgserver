dbuild:
	@docker build -t bimgserver .

build-linux:
	@go build -o bin/bimgserver-linux .

build: 
	@go build -o bin/bimgserver .

dev:
	@docker run -p 9080:90890 -it -v $(shell pwd):/go/src/app bimgserver bash

generate: 
	@protoc --twirp_out=paths=source_relative:. --go_out=paths=source_relative:. rpc/service.proto
