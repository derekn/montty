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

var (
	addr          string
	title         string
	version       string
	clients       *Clients
)

func readStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		fmt.Printf("%s\n", line)
		clients.Broadcast(append(line, '\n'))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error:", err)
	}
}

func init() {
	flag.StringVarP(&addr, "addr", "a", "127.0.0.1:8000", "ip:port to listen on")
	flag.StringVarP(&title, "title", "t", "", "app title")
	flag.BoolP("help", "h", false, "display usage help")
	flag.BoolP("version", "v", false, "display version")
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "pipe stdin to browser\n\n")
		fmt.Fprintf(os.Stderr, "usage: prog ... | montty [options...]\n\n")
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
			fmt.Printf("montty v%s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
			os.Exit(0)
		}
	}
	flag.Parse()

	clients = NewClients()
	registerRoutes()
	go readStdin()

	log.Printf("listening on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
