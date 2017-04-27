[![Build Status](https://travis-ci.org/cad/vehicle-tracker-api.svg?branch=master)](https://travis-ci.org/cad/vehicle-tracker-api)
# vehicle-tracker-api

## Get dependencies
`$ cd vehicle-tracker-api/`

`$ go get`

## Generate specs
`$ go generate`

## Build
Build for x86.

`$ go build -o vehicle-tracker`

Build and static link

`$ go build -ldflags "-linkmode external -extldflags -static" -o vehicle-tracker`

## Run

### On First Run
Create a user.

`$ ./vehicle-tracker createsuperuser --email example@example.com --password 1234`

### Start the api 
Run

`$ ./vehicle-tracker run`

