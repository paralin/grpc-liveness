# GRPC Liveness Checks

[![GoDoc Widget]][GoDoc] [![Go Report Card Widget]][Go Report Card]

[GoDoc]: https://godoc.org/github.com/paralin/grpc-liveness
[GoDoc Widget]: https://godoc.org/github.com/paralin/grpc-liveness?status.svg
[Go Report Card Widget]: https://goreportcard.com/badge/github.com/paralin/grpc-liveness
[Go Report Card]: https://goreportcard.com/report/github.com/paralin/grpc-liveness

## Introduction

**grpc-liveness** includes a GRPC service definition and executable to probe GRPC services for readiness and aliveness.

The intent is to use this in a Kubernetes exec check, to determine if a GRPC-only service is online.

## Getting Started

You can get the checker like so:

```
go get -u -v github.com/paralin/grpc-liveness/checker
```

To implement the checking service in your project, see the example directory.

