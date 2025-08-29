# go-json-rpc

https://pkg.go.dev/golang.org/x/exp/jsonrpc2 を使ってみるだけ


- `Serve` で RPC サーバーが作れる
  - この場合、クライアントにリクエストを送る事はできない
- `Dial` で双方向通信の接続が作れる
  - クライアントはこっちを使う
- ソケットやIOを開くには `NetListener`, `NetDialer`, を使う
  - `NetDialer` は `net.Pipe` を使った実装だけど、Accept しないと Dial できない事に注意
- `ConnectionOptions.Frame` で、Content-Type のありなしを選べる
- リクエストの受信は `ConnectionOptions.Handler` で行う


## Example

```
# terminal A
❯❯ go run . -mode tcp-server
2025/08/29 21:41:18 [tcp-server] got request ping
2025/08/29 21:41:18 [tcp-server] got request hello

# terminal B
❯❯ nc 0.0.0.0 1234 <<EOF
{"jsonrpc": "2.0", "method": "ping", "id": 1}
{"jsonrpc": "2.0", "method": "hello", "params": ["world"], "id": 1}
EOF
{"jsonrpc":"2.0","id":1,"result":"pong"}{"jsonrpc":"2.0","id":1,"result":"hello world"}
```

```
# terminal A
❯❯ go run . -mode tcp-server
2025/08/29 21:43:13 [tcp-server] got request ping
2025/08/29 21:43:27 [tcp-server] got request hello

# terminal B
❯❯ go run . -mode tcp-client
ping
2025/08/29 21:43:13 [tcp-client] command [ping]
2025/08/29 21:43:13 [tcp-client] got result: "pong"
pong
hello world
2025/08/29 21:43:27 [tcp-client] command [hello world]
2025/08/29 21:43:27 [tcp-client] got result: "hello world"
```

```
❯❯ go run . -mode stdio-server <<EOF
{"jsonrpc": "2.0", "method": "ping", "id": 1}
{"jsonrpc": "2.0", "method": "hello", "params": ["world"], "id": 1}
EOF
2025/08/29 21:45:49 [stdio-server] got request ping
{"jsonrpc":"2.0","id":1,"result":"pong"}
2025/08/29 21:45:49 [stdio-server] got request hello
{"jsonrpc":"2.0","id":1,"result":"hello world"}
```

```
❯❯ go run . -mode bidi-server
2025/08/29 21:53:13 [bidi-server] got request ping
2025/08/29 21:53:16 [bidi-server] got request hello
ping
2025/08/29 21:53:18 [bidi-server] command [ping]
2025/08/29 21:53:18 [bidi-server] got result: "pong"
hello world
2025/08/29 21:53:21 [bidi-server] command [hello world]
2025/08/29 21:53:21 [bidi-server] got result: "hello world"

❯❯ go run . -mode bidi-client
ping
2025/08/29 21:53:13 [bidi-client] command [ping]
2025/08/29 21:53:13 [bidi-client] got result: "pong"
hello world
2025/08/29 21:53:16 [bidi-client] command [hello world]
2025/08/29 21:53:16 [bidi-client] got result: "hello world"
2025/08/29 21:53:18 [bidi-client] got request ping
2025/08/29 21:53:21 [bidi-client] got request hello
```
