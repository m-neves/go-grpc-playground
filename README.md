# gRPC playground

Go gRPC playground for learning and testing

## Makefile Commands

Compile protobuf
`make all`

Start client
`client-start`

Start server
`server-start`

## Client Commands

### Type inside console

Unary request
`-cmd=unary -msg="place to rest"`

Server streaming
`-cmd=sstream -msg=place to rest`

Client streaming
`-cmd=cstream`

Bidirectional streaming
`-cmd=bidi`
