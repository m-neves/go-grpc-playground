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

#### Unary request

`-cmd=unary -msg="place to rest"`

#### Unary request with gRPC error

To raise an error:
`-cmd=errunary -msg=doerr`
To get a correct message
`-cmd=errunary -msg=correct`

#### Unary request with timeout

To timeout:
`-cmd=unarytimeout`
To get normal response:
`-cmd=unarytimeout -msg=ok`

#### Server streaming

`-cmd=sstream -msg=place to rest`

#### Client streaming

`-cmd=cstream`

#### Bidirectional streaming

`-cmd=bidi`

Server streaming
```-cmd=sstream -msg=place to rest```

Client streaming
```-cmd=cstream```

Bidirectional streaming
```-cmd=bidi```
