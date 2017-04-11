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
Run

`$ ./vehicle-tracker serve`

