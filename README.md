# load

> A toy implementation of a load testing tool.

Each worker creates it's own http client for all requests. This allows for all requests in that worker to reuse connections efficiently.

## Running

```bash
go run main.go -url http://localhost:8080 -workers 10 -duration 20s -verbose

go run main.go -help
```

```
  -duration duration
        Test duration (default 10s)
  -url string
        URL to load test (default "http://example.com")
  -verbose
        Verbose output
  -workers int
        Number of workers (default 1)
```

## Building

To build the `load` CLI tool, run the following command:

```sh
go build -o load main.go

# build for mac+win+linux in one go
./build.sh
```

This will generate an executable named `load_platform_arch` in the current directory.