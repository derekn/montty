# MonTTY

![GitHub Release](https://img.shields.io/github/v/release/derekn/montty)
![GitHub License](https://img.shields.io/github/license/derekn/montty)

Monitor stdout from CLI commands in browser.

## Installation

Download from Github [releases](https://github.com/derekn/montty/releases/latest).  
or, install using Go:

```shell
go install https://github.com/derekn/montty@latest
```

## Usage

```shell
command... | montty -a :8080

# example
i=1; while :; do echo "Hello, world $i"; ((i++)); done | montty
```

### Arguments

```shell
-a, --address string   ip:port to listen on (default "127.0.0.1:8000")
-t, --title string     app title
-b, --buffer int       history lines to buffer (default 500)
    --css string       custom CSS URL
-h, --help             display usage help
-v, --version          display version
```
