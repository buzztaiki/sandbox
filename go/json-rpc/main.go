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
	switch *mode {
	case "tcp-server":
		f = tcpServer
	case "stdio-server":
		f = stdioServer
	case "bidi-server":
		f = bidiServer
	case "tcp-client":
		f = tcpClient
	case "stdio-client":
		f = stdioClient
	case "bidi-client":
		f = bidiClient
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

func prefixLogger(prefix string) *log.Logger {
	return log.New(os.Stderr, "["+prefix+"] ", log.LstdFlags|log.Lmsgprefix)
}

func handler(logger *log.Logger) jsonrpc2.Handler {
	f := func(ctx context.Context, r *jsonrpc2.Request) (any, error) {
		logger.Println("got request", r.Method)
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
	}
	return jsonrpc2.HandlerFunc(f)
}

func logOnlyHandler(logger *log.Logger) jsonrpc2.Handler {
	f := func(ctx context.Context, r *jsonrpc2.Request) (any, error) {
		logger.Println("got request", r.Method)
		return nil, nil
	}
	return jsonrpc2.HandlerFunc(f)
}

func readCommandLoop(ctx context.Context, logger *log.Logger, r io.Reader, con *jsonrpc2.Connection) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		xs := strings.Split(line, " ")
		if len(xs) == 0 {
			continue
		}
		logger.Println("command", xs)

		var res json.RawMessage
		con.Call(ctx, xs[0], xs[1:]).Await(ctx, &res)
		logger.Println("got result:", string(res))

	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func tcpServer(ctx context.Context, framer jsonrpc2.Framer) error {
	logger := prefixLogger("tcp-server")

	l, err := jsonrpc2.NetListener(ctx, "tcp", ":1234", jsonrpc2.NetListenOptions{})
	if err != nil {
		return err
	}
	server, err := jsonrpc2.Serve(ctx, l, jsonrpc2.ConnectionOptions{Framer: framer, Handler: handler(logger)})
	if err != nil {
		return err
	}
	server.Wait()

	return nil
}

func stdioServer(ctx context.Context, framer jsonrpc2.Framer) error {
	logger := prefixLogger("stdio-server")

	l, err := jsonrpc2.NetPipe(ctx)
	if err != nil {
		return err
	}

	server, err := jsonrpc2.Serve(ctx, l, jsonrpc2.ConnectionOptions{
		Framer:  framer,
		Handler: handler(logger),
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
	logger := prefixLogger("tcp-client")

	d := jsonrpc2.NetDialer("tcp", "localhost:1234", net.Dialer{})
	con, err := jsonrpc2.Dial(ctx, d, jsonrpc2.ConnectionOptions{Framer: framer, Handler: logOnlyHandler(logger)})
	if err != nil {
		return err
	}
	defer con.Close()
	return readCommandLoop(ctx, logger, os.Stdin, con)
}

func stdioClient(ctx context.Context, framer jsonrpc2.Framer) error {
	logger := prefixLogger("tcp-client")

	l, err := jsonrpc2.NetPipe(ctx)
	if err != nil {
		return err
	}
	go func() {
		pipe, err := l.Accept(ctx)
		logger.Println("accepted pipe")
		if err != nil {
			logger.Println("error accepting pipe:", err)
		}
		go io.Copy(pipe, os.Stdin)
		go io.Copy(os.Stdout, pipe)
	}()

	d := l.Dialer()
	con, err := jsonrpc2.Dial(ctx, d, jsonrpc2.ConnectionOptions{Framer: framer, Handler: logOnlyHandler(logger)})
	if err != nil {
		return err
	}
	defer con.Close()

	nl, err := net.Listen("tcp", ":1234")
	if err != nil {
		return err
	}

	for {
		ncon, err := nl.Accept()
		logger.Print("accepted", ncon.RemoteAddr())
		if err != nil {
			return err
		}

		go readCommandLoop(ctx, logger, ncon, con)
	}
}

type rwcDialer struct {
	rwc io.ReadWriteCloser
}

func (d *rwcDialer) Dial(context.Context) (io.ReadWriteCloser, error) {
	return d.rwc, nil
}

func bidiServer(ctx context.Context, framer jsonrpc2.Framer) error {
	logger := prefixLogger("bidi-server")

	l, err := jsonrpc2.NetListener(ctx, "tcp", ":1234", jsonrpc2.NetListenOptions{})
	if err != nil {
		return nil
	}

	rwc, err := l.Accept(ctx)
	con, err := jsonrpc2.Dial(ctx, &rwcDialer{rwc}, jsonrpc2.ConnectionOptions{Framer: framer, Handler: handler(logger)})
	if err != nil {
		return err
	}
	defer con.Close()

	return readCommandLoop(ctx, logger, os.Stdin, con)
}

func bidiClient(ctx context.Context, framer jsonrpc2.Framer) error {
	logger := prefixLogger("bidi-client")

	d := jsonrpc2.NetDialer("tcp", "localhost:1234", net.Dialer{})
	con, err := jsonrpc2.Dial(ctx, d, jsonrpc2.ConnectionOptions{Framer: framer, Handler: handler(logger)})
	if err != nil {
		return err
	}
	defer con.Close()

	return readCommandLoop(ctx, logger, os.Stdin, con)
}
