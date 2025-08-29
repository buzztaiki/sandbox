package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/exp/jsonrpc2"
)

func main() {
	mode := flag.String("mode", "tcp-server", "tcp-server,stdio-server")
	useHeader := flag.Bool("use-header", false, "use header framer")
	flag.Parse()

	ctx := context.Background()
	var f func(context.Context, jsonrpc2.Framer) error
	switch (*mode) {
	case "tcp-server":
		f = tcpServer
	case "stdio-server":
		f = stdioServer
	case "tcp-client":
		f = tcpClient
	case "stdio-client":
		f = stdioClient
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
				var params []string
				if err := json.Unmarshal(r.Params, &params); err != nil {
					return nil, err
				}
				return "hello " + strings.Join(params, " "), nil
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

func stdioServer(ctx context.Context, framer jsonrpc2.Framer) error {
	l, err := jsonrpc2.NetPipe(ctx)
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

	pipe, err := l.Dialer().Dial(ctx)
	if err != nil {
		return err
	}
	go io.Copy(pipe, os.Stdin)
	go io.Copy(os.Stdout, pipe)

	server.Wait()

	return nil
}

func tcpClient(ctx context.Context, framer jsonrpc2.Framer) error {
	d := jsonrpc2.NetDialer("tcp", "localhost:1234", net.Dialer{})
	con, err := jsonrpc2.Dial(ctx, d, jsonrpc2.ConnectionOptions{
		Framer: framer,
		Handler: jsonrpc2.HandlerFunc(func(ctx context.Context, r *jsonrpc2.Request) (any, error) {
			log.Println("[tcpClient] got request", r.Method)
			return nil, nil
		}),
	})
	if err != nil {
		return err
	}
	defer con.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		xs := strings.Split(line, " ")
		if len(xs) == 0 {
			continue
		}
		var res json.RawMessage
		con.Call(ctx, xs[0], xs[1:]).Await(ctx, &res)
		log.Println("[tcpClient] got result:", string(res))

	}
	if scanner.Err() != nil {
		return scanner.Err()
	}



	return nil
}

func stdioClient(ctx context.Context, framer jsonrpc2.Framer) error {
	l, err := jsonrpc2.NetPipe(ctx)
	if err != nil {
		return err
	}
	go func() {
		pipe, err := l.Accept(ctx)
		log.Println("accepted pipe")
		if err != nil {
			log.Println("error accepting pipe:", err)
		}
		go io.Copy(pipe, os.Stdin)
		go io.Copy(os.Stdout, pipe)
	}()

	d := l.Dialer()
	con, err := jsonrpc2.Dial(ctx, d, jsonrpc2.ConnectionOptions{
		Framer: framer,
		Handler: jsonrpc2.HandlerFunc(func(ctx context.Context, r *jsonrpc2.Request) (any, error) {
			log.Println("[tcpClient] got request", r.Method)
			return nil, nil
		}),
	})
	if err != nil {
		return err
	}
	defer con.Close()


	nl, err := net.Listen("tcp", ":1234")
	if err != nil {
		return err
	}

	for {
		nc, err := nl.Accept()
		log.Print("accepted", nc.RemoteAddr())
		if err != nil {
			return err
		}

		go func() {
			scanner := bufio.NewScanner(nc)
			for scanner.Scan() {
				line := scanner.Text()
				xs := strings.Split(line, " ")
				if len(xs) == 0 {
					continue
				}
				log.Println("command", xs)

				var res json.RawMessage
				con.Call(ctx, xs[0], xs[1:]).Await(ctx, &res)
				log.Println("[tcpClient] got result:", string(res))

			}
			if scanner.Err() != nil {
				log.Println("scanner error:", scanner.Err())
			}
		}()
	}
}
