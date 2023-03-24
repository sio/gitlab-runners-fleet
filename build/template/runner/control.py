#!/usr/bin/env python3
'''
Simplified GitLab runner metrics
'''


METRICS_INTERNAL_URL = 'http://localhost:9252/metrics'
METRICS_NAMES = [
    'gitlab_runner_jobs',
    'gitlab_runner_jobs_total',
]
UNREGISTER_FIFO = '/run/unregister-runners.fifo'
UNREGISTER_COMMAND = 'UNREGISTER\n'
UNREGISTER_DELAY = 1  # seconds

import json
import re
import urllib.request
from http.server import HTTPServer, BaseHTTPRequestHandler
from socket import gethostname
from time import sleep


class ControlAPI(BaseHTTPRequestHandler):
    '''Translate metrics from local Prometheus exporter to public JSON'''
    def do_GET(self):
        if self.path != '/metrics':
            self.send_http(code=401, body='HTTP 401: Unauthorized')
            return
        try:
            self.send_http(
                body=json.dumps(get_metrics()),
                headers={'Content-Type': 'application/json'},
            )
        except Exception:
            self.send_http(code=500)

    def do_POST(self):
        if self.path != '/unregister' or self.headers.get('Host') != gethostname():
            self.send_http(code=401, body='HTTP 401: Unauthorized')
            return
        try:
            with open(UNREGISTER_FIFO, 'w') as fifo:
                fifo.write(UNREGISTER_COMMAND)
            sleep(UNREGISTER_DELAY)
            self.send_http('OK')
        except Exception:
            self.send_http(code=500)

    def send_http(self, body='', code=200, headers=None):
        if isinstance(body, str):
            body = body.encode()
        if headers is None:
            headers = dict()
        self.send_response(code)
        for header, value in headers.items():
            self.send_header(header, value)
        self.send_header('Content-Length', len(body))
        self.end_headers()
        self.wfile.write(body)


METRICS_REGEX = {name: re.compile(r'^%s{.*} (\d+)' % name) for name in METRICS_NAMES}
def get_metrics():
    response = urllib.request.urlopen(METRICS_INTERNAL_URL)
    if response.status != 200:
        raise ValueError(f'Internal metrics url returned HTTP {response.status}')

    metrics = {name: 0 for name in METRICS_NAMES}
    for line in response:
        for metric, regex in METRICS_REGEX.items():
            match = regex.match(line.decode())
            if match:
                metrics[metric] += int(match.group(1))
    return dict(metrics)


def main():
    address = ''
    port = 8080
    httpd = HTTPServer((address, port), ControlAPI)
    httpd.serve_forever()


if __name__ == '__main__':
    main()
