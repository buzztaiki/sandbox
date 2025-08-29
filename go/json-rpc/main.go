package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"maps"
	"net"
	"os"
	"slices"
	"strings"
	"time"

	"golang.org/x/exp/jsonrpc2"
)

func main() {
	modes := map[string]func(context.Context, jsonrpc2.Framer) error{
		"tcp-server":   tcpServer,
		"stdio-server": stdioServer,
		"bidi-server":  bidiServer,
		"tcp-client":   tcpClient,
		"bidi-client":  bidiClient,
	}

	mode := flag.String("mode", "tcp-server", strings.Join(slices.Sorted(maps.Keys(modes)), " "))
	useHeader := flag.Bool("use-header", false, "use header framer")
	flag.Parse()

	ctx := context.Background()
	f, ok := modes[*mode]
	if !ok {
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
			return "hello " + string(r.Params), nil
		}
		if r.Method == "sleep" {
			var params []int
			if err := json.Unmarshal(r.Params, &params); err != nil {
				return nil, err
			}
			time.Sleep(time.Duration(params[0]) * time.Second)
			return "awaiked", nil
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
		xs := strings.SplitN(line, " ", 2)
		if len(xs) == 0 {
			continue
		}

		method := xs[0]
		var params json.RawMessage
		if len(xs) > 1 {
			if err := json.Unmarshal([]byte(xs[1]), &params); err != nil {
				log.Println("invalid params:", err)
				continue
			}
		}

		logger.Println("send command", method, string(params))
		var res json.RawMessage
		if err := con.Call(ctx, method, params).Await(ctx, &res); err != nil {
			logger.Println("error:", err)
			continue
		}
		logger.Println("got result:", string(res))

	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func bindStream(ctx context.Context, l jsonrpc2.Listener, r io.Reader, w io.Writer) error {
	rwc, err := l.Accept(ctx)
	if err != nil {
		return err
	}
	go io.Copy(rwc, r)
	go io.Copy(w, rwc)
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
