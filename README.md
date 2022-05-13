# gRPC Web Go client
[![CircleCI](https://circleci.com/gh/ktr0731/grpc-web-go-client.svg?style=svg)](https://circleci.com/gh/ktr0731/grpc-web-go-client)  

gRPC Web client written in Go

## Usage
The server is [here](github.com/ktr0731/grpc-test).  

Send an unary request.

``` go
client := grpcweb.NewClient("localhost:50051")

in, out := new(api.SimpleRequest), new(api.SimpleResponse)
in.Name = "ktr"

// You can get the endpoint from grpcweb.ToEndpoint function with descriptors.
// However, I write directly in this example.
req := grpcweb.NewRequest("/api.Example/Unary", in, out)

res, err := client.Unary(context.Background(), req)
if err != nil {
  log.Fatal(err)
}

// hello, ktr
fmt.Println(res.Content.(*api.SimpleResponse).GetMessage())
```
