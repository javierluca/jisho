
.PHONY: build doc fmt lint dev test vet godep install bench

PKG_NAME=$(shell basename `pwd`)

install:
	go get -t -v ./...

fmt:
	go fmt ./...

godep:
	godep save ./...