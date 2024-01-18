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

You can find the pre-built binaries on the [Releases overview](https://github.com/anthonyringoet/load/releases).

## Output

Something like the following:

```
2xx requests: 426, non 2xx requests: 0
Total requests: 426
Avg requests/worker: 426

Min latency: 15.676958ms
Max latency: 93.211042ms
Median latency: 22.681084ms
90th percentile latency: 26.985ms
95th percentile latency: 28.549583ms
99th percentile latency: 34.739958ms
Average latency: 23.434747ms
```
