#!/bin/bash

build: binary

binary:
	@echo "building binary.."
	@go build -tags static_all .


clean:
	@echo "cleaning ..."
	@rm -f rll-url-shortener
	@rm -rf vendor
	@rm -f go.sum


install:
	@echo "Installing dependencies...."
	@rm -rf vendor
	@rm -f Gopkg.lock
	@rm -f glide.lock
	@go mod tidy && go mod download && go mod vendor

test:
	@go test $$(go list ./... | grep -v /vendor/) -cover

test-cover:
	@go test $$(go list ./... | grep -v /vendor/) -coverprofile=cover.out && go tool cover -html=cover.out ; rm -f cover.out

coverage:
	@go test -covermode=count -coverprofile=count.out fmt; rm -f count.out

start:
	@go run main.go serve-http

db.migrate.create:
	@go run main.go db:migrate create $(name) sql

db.migrate.up:
	@go run main.go db:migrate up

docker.compose.up:
	docker-compose -f deployment/docker-compose.yaml --project-directory . up -d --build

docker.compose.down:
	docker-compose -f deployment/docker-compose.yaml --project-directory . down 

app.serve.http:
	./rll-url-shortener serve-http