# TicketingSystem

A Sample Go program made for educational purposes. This relies on the local filesystem for storage/persistence so obviously this is not meant to be deployed to a production production-like environment.

### Prerequisites

You must have Go installed on your system. If not [see here](https://golang.org/doc/install)


### Installing

```
go get github.com/aaronireland/TicketingSystem
```

### Usage

Given a ticket request batch file (see below)
```
$GOPATH/bin/TicketingSystem -file /path/to/example.txt
```

Display more information about the batch request
```
$GOPATH/bin/TicketingSystem -file /path/to/example.txt -verbose
```

### Example File

example.txt

```
6 6
3 5 5 3
4 6 6 4
2 8 8 2
6 6

Smith 2
Jones 5
Davis 6
Wilson 100
Johnson 3
Williams 4
Brown 8
Miller 12
```

example_with_event_name.txt

```
6 6
3 5 5 3
4 6 6 4
2 8 8 2
6 6

Smith 2
Jones 5
Davis 6
Wilson 100
Johnson 3
Williams 4
Brown 8
Miller 12

event Example Event
```

Different theater layouts will persist to different files. Each event has its own data directory so the same theater layout can be use for different
reservation lists

## Running the tests

Note: TO-Do
```
go test
```
