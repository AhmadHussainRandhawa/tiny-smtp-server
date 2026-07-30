<div align="center">

# Tiny SMTP Server

**Receive email from raw SMTP clients.**

*A minimal SMTP server written in Go that implements the receiving side of the SMTP protocol—from `HELO` to `QUIT`—and stores accepted messages as `.eml` files.*

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white">
  <img alt="Protocol" src="https://img.shields.io/badge/Protocol-SMTP-2563EB">
  <img alt="Architecture" src="https://img.shields.io/badge/Architecture-State%20Machine-7C3AED">
  <img alt="Dependencies" src="https://img.shields.io/badge/Standard%20Library-Only-16A34A">
  <img alt="License" src="https://img.shields.io/badge/License-MIT-F59E0B">
</p>

</div>

---

## 🎬 See it in action

https://github.com/user-attachments/assets/2f312d58-e20c-42a4-a8e1-62d4f4005311

*A live session: a client connects, walks through `EHLO` → `MAIL FROM` → `RCPT TO` → `DATA` → `QUIT`, and the message lands on disk as a `.eml` file — plus what happens when a client sends things out of order.*

--- 

## Why This Project Exists

When an email reaches its destination, a receiving mail server takes over.

Its job is to speak the **Simple Mail Transfer Protocol (SMTP)**, validate the client's commands, receive the message, and decide whether to accept or reject it.

This project implements that receiving workflow by:

* Accepting SMTP connections over TCP
* Enforcing the correct SMTP command sequence
* Receiving message data
* Persisting accepted messages as `.eml` files

It's a minimal implementation of the server-side protocol that powers email delivery.

---

## Getting Started

### 1. Clone the repository

```bash
git clone git@github.com:AhmadHussainRandhawa/tiny-smtp-server.git
cd tiny-smtp-server
```

### 2. Start the server

```bash
go run .
```

The server listens on:

```text
127.0.0.1:2525
```

Accepted messages are stored in the `messages/` directory as `.eml` files.

### 3. Connect with an SMTP client

Open another terminal and connect to the server:

```bash
nc 127.0.0.1 2525
```

Then start an SMTP session:

```text
EHLO localhost
MAIL FROM:<alice@example.com>
RCPT TO:<bob@example.com>
DATA
Subject: Hello

This message was received by the Tiny SMTP Server.
.
QUIT
```

After the terminating `.`, the server accepts the message and saves it as a new `.eml` file inside the `messages/` directory.

---

## The state machine

This is the actual engineering of the project: what's a legal next command, and what gets rejected.

```
                     ┌───────────────┐
                     │   CONNECTED   │
                     └───────┬───────┘
                             │ EHLO / HELO <domain>
                             ▼
                     ┌───────────────┐
        ┌───────────▶│    GREETED    │
        │            └───────┬───────┘
        │                    │ MAIL FROM:<addr>
        │                    ▼
        │            ┌───────────────┐
        │            │  SENDER SET   │
        │            └───────┬───────┘
        │                    │ RCPT TO:<addr>   (repeatable)
        │                    ▼
        │            ┌───────────────┐
        │            │ RECIPIENT(S)  │
        │            │     SET       │
        │            └───────┬───────┘
        │                    │ DATA
        │                    ▼
        │            ┌───────────────┐
        │            │  RECEIVING    │
        │            │    BODY       │
        │            └───────┬───────┘
        │                    │ <CRLF>.<CRLF>
        │                    │ → write messages/*.eml
        └────────────────────┘
             session resets — ready for the next MAIL FROM

  Out-of-order command  →  503
  Malformed syntax       →  501
  Unrecognized command   →  500
  QUIT (from any state)  →  221, connection closes
```

---

## 💬 Example SMTP Session

A successful SMTP transaction between a client and the server looks like this:

```text
Client                                    Tiny SMTP Server
──────────────────────────────────────────────────────────────────────

                                          220 localhost Tiny SMTP Server Ready

EHLO localhost                  ───────▶
                                ◀─────── 250 localhost Hello

MAIL FROM:<farman@example.com>   ───────▶
                                ◀─────── 250 Sender accepted

RCPT TO:<ahmad@example.com>       ───────▶
                                ◀─────── 250 Recipient accepted

DATA                            ───────▶
                                ◀─────── 354 End data with <CRLF>.<CRLF>

From: farman <farman@example.com>
To: ahmad <ahmad@example.com>
Subject: Important Discussion


Aslam U Alaikum!

How are you? i have a meetup plan.

Regards,
Ahmad
.
                                ───────▶
                                ◀─────── 250 Message accepted

QUIT                            ───────▶
                                ◀─────── 221 Bye
```

After the terminating `.` is received, the server writes the message to the `messages/` directory as an `.eml` file and resets the session, ready to receive the next message.

---

## Design Decisions

This project intentionally focuses on implementing the **core SMTP protocol** rather than building a production-ready mail server.

| Decision                            | Rationale                                                                                                           |
| ----------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| **Single client connection**        | Keeps the implementation focused on the SMTP protocol instead of concurrency.                                       |
| **No authentication or TLS**        | Authentication and transport security are separate concerns from the SMTP state machine explored here.              |
| **Localhost only**                  | Designed for learning, experimentation, and protocol exploration in a controlled environment.                       |
| **Messages stored as `.eml` files** | Makes it easy to inspect exactly what the server receives without introducing mailbox storage complexity.           |
| **Standard library only**           | Demonstrates that a functional SMTP server can be built using Go's built-in networking and file I/O packages alone. |

---

## What's Intentionally Out of Scope

This project is designed to explore the SMTP protocol itself—not to replace a production mail server.

| Not implemented                                | Why                                                                                          |
| ---------------------------------------------- | -------------------------------------------------------------------------------------------- |
| **Concurrent client handling**                 | Keeps the implementation focused on the SMTP protocol rather than connection management.     |
| **SMTP extensions (`STARTTLS`, `AUTH`, etc.)** | The goal is to implement the core SMTP workflow before protocol extensions.                  |
| **Mailbox delivery (MDA)**                     | Messages are stored as `.eml` files instead of being delivered to user mailboxes.            |
| **Message relay**                              | The server accepts and stores mail locally; it does not forward messages to other MTAs.      |
| **Spam filtering or security policies**        | SPF, DKIM, DMARC, greylisting, and reputation systems are outside the scope of this project. |
| **RFC-complete SMTP implementation**           | Only the core command set required to understand SMTP message reception is implemented.      |

---

## Part of a series

This project is the final installment in a three-project journey to understand how email works—from sending a message, to discovering where it should go, to receiving and storing it.

| Project | What it proves | Status |
|---|---|:---:|
| **[tiny-smtp-client](https://github.com/AhmadHussainRandhawa/tiny-smtp-client)** | Sending mail, from a raw socket up | ✅ Complete |
| **[mx-lookup-tool](https://github.com/AhmadHussainRandhawa/mx-lookup-tool)** | DNS resolution of mail routing (MX + A records) | ✅ Complete |
| **[tiny-smtp-server](https://github.com/AhmadHussainRandhawa/tiny-smtp-server)** *(this repo)* | Receiving mail — including malformed input | ✅ Complete |

---

<div align="center">

**License:** [MIT](./LICENSE)

</div>

