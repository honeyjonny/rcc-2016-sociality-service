
# rcc-2016-sociality-service

Simple application that emulates social network behaviour

with three vulnerabilities :)


## Application use:

	postgress = 9.4v


### Golang libraries:

	gin
	gorm
	github.com/lib/pq


### Install:

1. Clone this repository in ```$GOPATH/src/```
2. Install dependency
    ``` go get -u "github.com/gin-gonic/gin" ```
    ``` go get -u "github.com/jinzhu/gorm" ```
    ``` go get -u "github.com/lib/pq" ```
3. Build / Run
    * Build
    ``` go build server.go ```
    ```./server ```
    * or run
    ``` go run server.go ```
