package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"

	"golang.org/x/exp/jsonrpc2"
)

func main() {
	mode := flag.String("mode", "tcp-server", "mode")
	useHeader := flag.Bool("use-header", false, "use header framer")
	flag.Parse()

	ctx := context.Background()
	var f func(context.Context, jsonrpc2.Framer) error
	switch (*mode) {
	case "tcp-server":
		f = tcpServer
	default:
		log.Fatal("unknown mode", *mode)
	}

	framer := jsonrpc2.RawFramer()
	if *useHeader {
		framer = jsonrpc2.HeaderFramer()
	}

	if err := f(ctx, framer); err != nil {
		log.Fatal(err)
	}
}

func tcpServer(ctx context.Context, framer jsonrpc2.Framer) error {
	l, err := jsonrpc2.NetListener(ctx, "tcp", ":1234", jsonrpc2.NetListenOptions{})
	if err != nil {
		return err
	}
	server, err := jsonrpc2.Serve(ctx, l, jsonrpc2.ConnectionOptions{
		Framer: framer,
		Handler: jsonrpc2.HandlerFunc(func(ctx context.Context, r *jsonrpc2.Request) (any, error) {
			if r.Method == "ping" {
				return "pong", nil
			}
			if r.Method == "hello" {
				var params struct {
					Name string `json:"name"`
				}
				if err := json.Unmarshal(r.Params, &params); err != nil {
					return nil, err
				}
				return "hello " + params.Name, nil
			}
			return nil, jsonrpc2.ErrMethodNotFound
		}),
	})
	if err != nil {
		return err
	}
	server.Wait()

	return nil
}
