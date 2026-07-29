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

	var greeted bool
	var mailFrom string
	var recipients []string
	var messageLines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", err)
			return
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			continue
		}

		fmt.Println("Client:", line)

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "EHLO", "HELO":
			if len(parts) < 2 {
				fmt.Fprint(conn, "501 HELO/EHLO requires a domain\r\n")
				continue
			}

			greeted = true

			fmt.Fprintf(conn, "250 %s Hello\r\n", parts[1])

		case "MAIL":
			if !greeted {
				fmt.Fprint(conn, "503 send HELO/EHLO first\r\n")
				continue
			}

			if !strings.HasPrefix(strings.ToUpper(line), "MAIL FROM:") {
				fmt.Fprint(conn, "501 Syntax: MAIL FROM:<address>\r\n")
				continue
			}

			mailFrom = strings.TrimSpace(line[len("MAIL FROM:"):])

			if mailFrom == "" {
				fmt.Fprint(conn, "501 Sender address is required\r\n")
				continue
			}

			recipients = nil

			fmt.Fprint(conn, "250 Sender accepted\r\n")

		case "RCPT":
			if mailFrom == "" {
				fmt.Fprint(conn, "503 send MAIL FROM first\r\n")
				continue
			}

			if !strings.HasPrefix(strings.ToUpper(line), "RCPT TO:") {
				fmt.Fprint(conn, "501 Syntax: RCPT TO:<address>\r\n")
				continue
			}

			recipient := strings.TrimSpace(line[len("RCPT TO:"):])

			if recipient == "" {
				fmt.Fprint(conn, "501 Recipient address is required\r\n")
				continue
			}

			recipients = append(recipients, recipient)

			fmt.Fprint(conn, "250 Recipient accepted\r\n")

		case "DATA":
			if mailFrom == "" {
				fmt.Fprint(conn, "503 Send MAIL FROM first\r\n")
				continue
			}

			if len(recipients) == 0 {
				fmt.Fprint(conn, "503 Send RCPT TO first\r\n")
				continue
			}

			fmt.Fprint(conn, "354 End data with <CRLF>.<CRLF>\r\n")

			messageLines = nil

			for {
				messageLine, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Failed to read message:", err)
					return
				}

				messageLine = strings.TrimRight(messageLine, "\r\n")

				if messageLine == "." {
					break
				}

				messageLines = append(messageLines, messageLine)
			}

			fmt.Println("\nMessage received:")
			fmt.Println(strings.Join(messageLines, "\n"))

			fmt.Fprint(conn, "250 Message accepted\r\n")

		default:
			fmt.Fprint(conn, "500 command not recognized\r\n")
		}

	}
}
