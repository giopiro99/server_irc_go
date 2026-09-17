# Go TCP IRC Server

A lightweight, asynchronous TCP chat server written in Go. This project implements core IRC-like mechanics to demonstrate advanced concurrency patterns, network I/O handling, and proper goroutine lifecycle management.

## Motivation
This project was built to explore Go's concurrency model in a networking context. Instead of relying on traditional memory sharing with mutexes (`sync.Mutex`), the architecture strictly adheres to Go's proverb: *"Do not communicate by sharing memory; instead, share memory by communicating."* 

All state management and client broadcasting are handled by a single, isolated Broadcaster goroutine using channels and the `select` statement, completely eliminating the risk of data races.

## Key Features
* **Mutex-Free Concurrency:** Centralized state management using a dedicated Broadcaster goroutine and channels.
* **Graceful Shutdown:** Intercepts system signals (`SIGINT`, `SIGTERM`) to actively close the main network listener, disconnect all active clients cleanly, and prevent goroutine leaks before exiting.
* **Request-Response Channel Pattern:** Validates unique client nicknames by passing dedicated response channels over a request channel, allowing safe, lock-free queries against the server state.
* **Environment Configuration:** Uses `.env` files to manage server properties (IP, Port, Password) securely.
* **Connection Handling:** Limits authentication attempts and safely manages client disconnects (EOF) without crashing the server.

## Getting Started

### Prerequisites
* Go 1.21 or higher
* Netcat (`nc`) for testing

### Installation & Setup
1. Clone the repository:
   ```bash
   git clone <your-repo-url>
   cd <repository-name>

    Set up the environment variables by copying the example file:

    cp .env.example .env

    (You can edit the .env file to change the default IP, PORT, or PASSWORD).

    Run the server:
   ❯ go run .

   If you want to check for data races or similar issues, use:
   ❯ go run -race .


Usage

Once the server is running, you can connect to it using multiple terminal windows via Netcat.

Connect a client:
nc 127.0.0.1 8190

    Insert the server password defined in the .env file.

    Choose a unique nickname (1 to 9 characters).

    Start typing to broadcast messages to all other connected users.

    Press Ctrl+C in the server terminal to observe the graceful shutdown process disconnecting all active clients.
