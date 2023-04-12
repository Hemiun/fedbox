#!/bin/bash

./bin/fedboxctl bootstrap reset

# Call twice! Important for fs storage
./bin/fedboxctl bootstrap
./bin/fedboxctl bootstrap

./bin/fedboxctl oauth client add -redirectUri=http://localhost:4000/ping

# admin actor
# ./bin/fedboxctl pub actor add -attributedTo=http://localhost:30000/actors/845cc562-8719-4477-bc00-d8a8e7d08a1b

# build in debug mode
go build -a -ldflags '-X main.version=feature-user-08ffae7' -gcflags="all=-N -l" -tags "dev storage_all" -o bin/fedbox ./cmd/fedbox/main.go
go build -a -ldflags '-X main.version=feature-user-08ffae7' -gcflags="all=-N -l" -tags "dev storage_all" -o bin/fedboxctl ./cmd/control/main.go

# regular build
#make download all

# debugger
#dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec bin/fedbox
#dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec bin/fedboxctl -- pub actor add  -attributedTo=http://localhost:30000/actors/7d5f07dd-1ebd-44db-adc5-5b2edc0011af
