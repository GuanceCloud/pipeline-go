# DataKit Web Pipeline Examples

Examples from `/usr/local/datakit/pipeline` for web and application-server logs.

## Nginx Access

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/nginx.p \
  --message '192.168.1.10 - john [10/Oct/2023:13:55:36 +0800] "GET /api/v1/users?id=42 HTTP/1.1" 200 1234 "https://example.com/" "Mozilla/5.0"' \
  --require-key client_ip \
  --require-key status_code \
  --expect status=OK
```

Key behavior: combined access logs extract HTTP fields, parse user-agent details, cast `status_code` and `bytes`, and map 2xx status to `OK`.

## Apache Access

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/apache.p \
  --message '10.0.0.1 - - [10/Oct/2023:13:55:36 +0800] "POST /login HTTP/1.1" 404 512' \
  --require-key ip_or_host \
  --require-key http_code \
  --expect http_code=404
```

Key behavior: extracts `ip_or_host`, request fields, and casts `http_code` to int.

## Tomcat JUL

```sh
./bin/pipeline-check \
  --script /usr/local/datakit/pipeline/tomcat.p \
  --message '19-May-2026 11:42:31.123 ERROR [main] org.apache.catalina.startup.Catalina Server startup failed' \
  --require-key status \
  --require-key thread_name \
  --require-key report_source \
  --expect status=ERROR
```

Key behavior: JUL application logs keep `status=ERROR`; access logs in the same script use `group_between(status_code, ...)` to derive normalized status.
