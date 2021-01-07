# wasm server

an HTTP server for wasm

## installation
- download zip

## usage
````
Usage of wasmserve
  -allow-origin string
        Allow specified origin (or * for all origins) to make requests to this server
  -http string
        HTTP bind address to serve (default ":8080")
  -tags string
        Build tags
````

## example
```cassandraql
cd wasmserver
go build -o wasmserver
./wasmserver -tags=example
```
open http://localhost:8080/ on your browser 

## inspired
https://github.com/hajimehoshi/wasmserve