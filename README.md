# Access verification service for CH clusters

The service enables you to proxy requests to HTTP endpoints of ClickHouse clusters
with advanced security settings.

Access is granted based on the following logic: 'Are requests from the
corresponding subnet allowed for the corresponding clickhouse cluster?'"

## Examples

Example of nginx config:

```nginx
# listening addr of the ch-proxy-auth
upstream ch_auth {
    server 127.0.0.1:8081;
}

server {
    ...
    server_name *.proxy.example.com;
    ...
    ...
    # tries to get clickhouse hostname from the Host header field
    set $ch_host "";
    if ($http_host ~* "^(ch-[a-z0-9]+)\.proxy\.example.com$") {
        set $ch_host "$1";
    }
    ...
    ...
    location = /auth {
        internal;
        proxy_pass              http://ch_auth;
        proxy_pass_request_body off;
        proxy_set_header        Content-Length "";
        # send real-ip and ch hostname to the ch-auth-proxy-svc
        proxy_set_header        X-Real-IP   $remote_addr;
        proxy_set_header        X-Server    $ch_host;
    }
    ......
    ......
    location / {
        auth_request /auth;
        proxy_pass http://$ch_host:8123;
        ......
    }
}
```
