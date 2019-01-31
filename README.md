# tcp_server_in_go

tcp server listen on 7788 - tcp server receive client query
```
=================================
HELLO TCP server
send 'quit' for close connection
=================================
> hello
received:hello
> quit
bye~%
```

http server listen on 7799 - get server status
```
========Client INFO========
client: 1
client addr: 127.0.0.1:65091
client connect time: 2019-01-31T20:36:04+08:00
client connect duration: 14s
processed requests: 1
request rate: 0.07/s

client: 2
client addr: 127.0.0.1:65121
client connect time: 2019-01-31T20:36:11+08:00
client connect duration: 7s
processed requests: 2
request rate: 0.27/s

========Server INFO========
current connections: 2
current remaining jobs: 0
total processed jobs: 3
```

## BUILD
```
make go-build
```

## RUN

```
make go-run
```