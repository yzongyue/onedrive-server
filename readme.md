### a simple server

- recv onedrive webhook
- communicate with clients over websocket

#### caddy config
```
example.com
{
    proxy / http://localhost:6500
    proxy /ws http://localhost:6500 {
        websocket
        transparent
    }
}

```

#### usage
go 1.12 or newer

```
go build . && onedrive-server
```
