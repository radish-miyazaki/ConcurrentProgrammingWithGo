package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

var r, _ = regexp.Compile("GET (.+) HTTP/1.1\r\n")

func main() {
	incommingConns := make(chan net.Conn)
	StartHTTPWorkers(3, incommingConns)

	srv, _ := net.Listen("tcp", "localhost:8080")
	defer srv.Close()

	for {
		conn, _ := srv.Accept()
		select {
		case incommingConns <- conn:
		default:
			fmt.Println("Server is busy")
			conn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n\r\n<html>Busy</html>"))
			conn.Close()
		}
	}
}

func handleHttpRequest(conn net.Conn) {
	buff := make([]byte, 1024)
	size, _ := conn.Read(buff)

	if r.Match(buff[:size]) {
		file, err := os.ReadFile(fmt.Sprintf("../resources/%s",
			r.FindSubmatch(buff[:size])[1],
		))

		if err == nil {
			fmt.Println("200")
			conn.Write(fmt.Appendf(nil, "HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n", len(file)))
			conn.Write(file)
		} else {
			fmt.Println("404")
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n<html>Not Found</html>"))
		}
	} else {
		fmt.Println("500")
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}

	conn.Close()
}

func StartHTTPWorkers(n int, incommingConns <-chan net.Conn) {
	for range n {
		go func() {
			for c := range incommingConns {
				handleHttpRequest(c)
			}
		}()
	}
}
