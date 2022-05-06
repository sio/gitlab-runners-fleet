#!/usr/bin/env python3
'''
Simplified GitLab runner metrics
'''


METRICS_INTERNAL_URL = 'http://localhost:9252/metrics'
METRICS_NAMES = [
    'gitlab_runner_jobs',
    'gitlab_runner_jobs_total',
]

import json
import re
import urllib.request
from http.server import HTTPServer, BaseHTTPRequestHandler


class MetricsHandler(BaseHTTPRequestHandler):
    '''Translate metrics from local Prometheus exporter to public JSON'''
    def do_GET(self):
        if self.path != '/metrics':
            self.send_response(401)
            self.end_headers()
            self.wfile.write(b'')
            return
        try:
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(get_metrics()).encode())
        except Exception:
            self.send_response(500)
            self.end_headers()
            self.wfile.write(b'')


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
    httpd = HTTPServer((address, port), MetricsHandler)
    httpd.serve_forever()


if __name__ == '__main__':
    main()
