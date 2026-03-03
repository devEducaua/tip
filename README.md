# tip protocol

a simple tcp/ip protocol, focused on simplicity. something in between nex and spartan.
but with no virtual hosting, and no domains in general.

*this is just a hobby project, and is still in development.

# spec

- port:     :1979
- scheme:   tip://
- encoding: utf-8
- tls:      none

# url scheme
```
scheme + ip+port + path
```

no domains.

## requests
```
ip:port path
```

## responses
```
code\ncontent
```

### status code
- 0 -> success
- 1 -> redirect
- 2 -> error


# todo
- mimetypes
- redirects
