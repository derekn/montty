package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	flag "github.com/spf13/pflag"
)

const appName = "montty"

var (
	addr          string
	title         string
	version       string
	clients       *Clients
	logBuffer     *LogBuffer
	logBufferSize int
)

func readStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		fmt.Printf("%s\n", line)
		clients.Broadcast(append(line, '\n'))
		logBuffer.AddLine(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error reading stdin:", err)
	}
}

func init() {
	flag.StringVarP(&addr, "addr", "a", "127.0.0.1:8000", "ip:port to listen on")
	flag.StringVarP(&title, "title", "t", "", "app title")
	flag.IntVarP(&logBufferSize, "buffer", "b", 500, "history lines to buffer")
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
	logBuffer = NewLogBuffer(int(logBufferSize))
	registerRoutes()
	go readStdin()

	log.Printf("listening on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
