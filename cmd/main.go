package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/buildkite/terminal-to-html"
	flag "github.com/spf13/pflag"
)

const appName = "montty"

type Args struct {
	Address       string
	Title         string
	LogBufferSize int
	CSSUrl        string
}

var (
	args      Args
	version   string
	clients   *Clients
	logBuffer *LogBuffer
)

func fmtOutput(output []byte) []byte {
	return append(terminal.Render(output), '\n')
}

func readStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		logBuffer.AddLine(line)
		fmt.Printf("%s\n", line)
		clients.Broadcast(fmtOutput(line))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error reading stdin:", err)
	}
}

func init() {
	flag.StringVarP(&args.Address, "address", "a", "127.0.0.1:8000", "ip:port to listen on")
	flag.StringVarP(&args.Title, "title", "t", "", "app title")
	flag.IntVarP(&args.LogBufferSize, "buffer", "b", 500, "history lines to buffer")
	flag.StringVar(&args.CSSUrl, "css", "", "custom CSS URL")
	flag.BoolP("help", "h", false, "display usage help")
	flag.BoolP("version", "v", false, "display version")
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "pipe stdin to browser\n\n")
		fmt.Fprintf(os.Stderr, "usage: prog ... | %s [options...]\n\n", appName)
		flag.PrintDefaults()
	}
	if version == "" {
		version = time.Now().Format("2006.1.2-dev")
	}
}

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			flag.Usage()
			os.Exit(0)
		case "--version", "-v":
			fmt.Printf("%s v%s %s/%s\n", appName, version, runtime.GOOS, runtime.GOARCH)
			os.Exit(0)
		}
	}
	flag.Parse()

	clients = NewClients()
	logBuffer = NewLogBuffer(int(args.LogBufferSize))
	registerRoutes()
	go readStdin()

	log.Printf("listening on http://%s\n", args.Address)
	log.Fatal(http.ListenAndServe(args.Address, nil))
}
