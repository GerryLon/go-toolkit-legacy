package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// send or receive from host:port
// Usage: socket -r/-s -h|--host host -p|--port port
func main() {
	opts, err := parseCmdLine()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if opts.mode == "r" {
		err = receive(&opts)
	} else if opts.mode == "s" {
		err = send(&opts)
	}

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

type Options struct {
	mode string // send or receive: s|r
	host string
	port int
}

func usage() string {
	return "Usage: /path/to/socket -r/-s -h host -p port"
}

func parseCmdLine() (Options, error) {
	opts := Options{}
	mode := "s" // default mode is send data

	r := flag.Bool("r", false, "use receive mode")
	s := flag.Bool("s", false, "use send mode")
	h := flag.String("h", "", "host: you want to connect")
	p := flag.Int("p", -1, "port: tcp port")
	// help := flag.Bool("help", false, "--help")
	flag.Parse()

	// set both r and s option, error
	if *r && *s {
		return opts, errors.New("can not set both r and s mode")
	}

	// no options or --help, return help
	if !*r && !*s && *p == -1 && *h == "" {
		return opts, errors.New(usage())
	}

	if *r {
		mode = "r"
	} else if *s {
		mode = "s"
	} else {
		mode = "s"
	}

	// check host
	if strings.TrimSpace(*h) == "" {
		*h = ":"
	}

	// check port
	if *p < 0 || *p > 65535 {
		return opts, errors.New("port should in [0, 65535]")
	}

	opts.mode = mode
	opts.host = *h
	opts.port = *p

	return opts, nil
}

func buildAddress(opts *Options) string {
	h := opts.host
	p := opts.port

	return fmt.Sprintf("%s:%d", h, p)
}

// you can test send by
// ./socket.out -s -p 7 -h 192.168.137.166, 7 is the echo service port
func send(opts *Options) error {
	addr, err := net.ResolveTCPAddr("tcp", buildAddress(opts))

	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	// conn.SetWriteBuffer(0) // will cause segment fault
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("receive End")
				return
			default:

				// read from conn
				s, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					fmt.Printf("< Err: %v\n", err)
				}
				fmt.Printf("< %s", s)
				fmt.Printf("> ")
			}
		}
	}(ctx)

	// read from stdin and send it to the host:port
	inputReader := bufio.NewReader(os.Stdin)
	var line string

	// send data to server, stop when get EOF
	fmt.Printf("> ")
	for {

		line, err = inputReader.ReadString('\n')

		if err == io.EOF {
			fmt.Println("End.")
			return err
		}

		if err != nil {
			fmt.Printf("Err: %v\n", err)
			continue
		}

		_, err := fmt.Fprintf(conn, "%s", line)

		if err != nil {
			fmt.Printf("Err: %v\n", err)
			continue
		}
	}

	return nil
}

// TODO
func receive(opts *Options) error {
	addr, err := net.ResolveTCPAddr("tcp", buildAddress(opts))

	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	defer listener.Close()

	if err != nil {
		return err
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
}
