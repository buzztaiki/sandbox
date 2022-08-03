#!/bin/python3

from http.server import BaseHTTPRequestHandler, HTTPServer

class Handler(BaseHTTPRequestHandler):
    def writeln(self, s):
        self.wfile.write(s.encode("utf-8"))
        self.wfile.write(b"\n")

    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-Type', 'text/plain; charset=utf-8')
        self.end_headers()

        self.writeln(f"path: {self.path}")
        self.writeln("headers:")
        for x in self.headers.items():
            self.writeln(f"  {x}")

with HTTPServer(("0.0.0.0", 8000), Handler) as server:
    server.serve_forever()
