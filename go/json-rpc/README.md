# go-json-rpc

https://pkg.go.dev/golang.org/x/exp/jsonrpc2 を使ってみるだけ

```
❯❯ nc 0.0.0.0 1234 <<EOF
∙ {"jsonrpc": "2.0", "method": "ping", "id": 1}
{"jsonrpc": "2.0", "method": "hello", "params": {"name":"world"}, "id": 1}
EOF
{"jsonrpc":"2.0","id":1,"result":"pong"}
{"jsonrpc":"2.0","id":1,"result":"hello world"}
```
