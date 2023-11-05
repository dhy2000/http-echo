# HTTP Echo

A tiny http server that always simply echos your request, written in Go.

## Build

```bash
go build
```

## Usage

Start an http server (default listening to address `0.0.0.0` and port `8000`) :

```bash
./http-echo
```

Start an http server listening to local loopback address `127.0.0.1` and port `8081`:

```bash
./http-echo -a 127.0.0.1:8081
```

Start an http server with name `YourAwesomeName`:

```bash
./http-echo -n YourAwesomeName
```

Start an https server listening to `0.0.0.0:4443` with certificate file `cert.pem` and key file `key.pem`:

```bash
./http-echo -a 0.0.0.0:4443 -c cert.pem -k key.pem
```

Print usage:

```bash
./http-echo -h
```