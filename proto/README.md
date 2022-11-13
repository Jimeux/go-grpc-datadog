# proto

### Protobuf definitions for servers.

## Build

`buf build`

## Generate Code

`buf generate`

## After adding a dependency

`buf mod update`

## Example usage of generated Go code

- `svc/second/go.mod`
    ```
    replace github.com/Jimeux/go-grpc-datadog/proto/go => ../../proto/go
    ```

- `svc/second/rpc/service.go`
    ```go
    import "github.com/Jimeux/go-grpc-datadog/proto/go/pb/second/v1"
    ```
