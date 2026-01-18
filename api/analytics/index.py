from http.server import BaseHTTPRequestHandler
import json
import numpy
import polars
import psycopg2

class handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps({'status': 'Would this work'}).encode())