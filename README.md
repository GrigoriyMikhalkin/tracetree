## Trace tree builder

Builds trees for traceId from logs

## How to run

First, build binary `go build`.

Binary can be run without arguments:

```
./tracetree
```

Then it reads input from stdin and writes back to stdout.

You can specify input file via flag `-infile=filepath`, and output file via `-outfile=filepath`, for example:

```
./tracetree -infile=filepath -outfile=filepath
```

# Run tests

```
go test
```
