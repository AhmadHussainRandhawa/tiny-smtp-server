package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:2525")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}

	defer listener.Close()

	err = os.MkdirAll("messages", 0755)
	if err != nil {
		fmt.Println("Failed to create messages directory:", err)
	}

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

			message := strings.Join(messageLines, "\r\n")
			fileName := filepath.Join(
				"messages", fmt.Sprintf("message-%s.eml", time.Now().Format("20060102-150405.000000000")))

			err := os.WriteFile(fileName, []byte(message), 0644)
			if err != nil {
				fmt.Println("Failed to save message:", err)
				fmt.Fprint(conn, "451 Failed to store message\r\n")
				continue
			}

			fmt.Println("Message saved to:", fileName)
			fmt.Fprint(conn, "250 Message accepted\r\n")

		default:
			fmt.Fprint(conn, "500 command not recognized\r\n")
		}

	}
}
