package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:2525")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}

	defer listener.Close()

	fmt.Println("SMTP server listening on 127.0.0.1:2525")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Failed to accept connection:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	fmt.Fprint(conn, "220 localhost Tiny SMTP Server Ready\r\n")

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", err)
			return
		}

		line = strings.TrimRight(line, "\r\n")

		fmt.Println("Client:", line)
	}
}
