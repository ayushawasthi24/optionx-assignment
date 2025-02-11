# Go WebSocket Server - OptionX Assignment

A server implementation using WebSocket protocol in Go. This server supports both broadcast and private messaging capabilities, with client management and ping/pong messages.

## Features

### Core Functionality

- Real-time bidirectional communication using WebSocket protocol
- Support for both broadcast and private messages
- Automatic client list synchronization
- Built-in ping/pong mechanism for connection health monitoring
- Graceful connection/disconnection handling and cleanup

### Message Types

- Broadcast messages to all connected clients
- Private messages between specific users
- Welcome messages for new connections
- System notifications for user joins/leaves
- Client list updates

### Security & Reliability

- Automatic connection timeout handling
- Proper connection cleanup on disconnection
- JSON message validation
- Username validation and uniqueness check

## Prerequisites

- Go 1.19 or higher
- `github.com/gorilla/websocket` package
- `github.com/brianvoe/gofakeit/v6` package for generating random usernames
- `github.com/google/uuid`

## Installation

1. Clone the repository:

```bash
git clone https://github.com/ayushawasthi24/optionx-assignment
cd optionx-assignment
```

2. Install dependencies:

```bash
go mod tidy
```

## Running the Server

1. Start the server:

```bash
go run main.go
```

2. The server will start on the configured port (default: 8080)

#### If using Docker

```bash
docker build -t optionx-assignment .
```

```bash
docker run -p 8080:8080 optionx-assignment
```

## API Reference

### WebSocket Endpoint

```
ws://localhost:8080/ws
```

### Message Formats

#### Regular Message

```json
{
  "type": "private",
  "sender": "your username",
  "receiver": "recipents client id",
  "content": "Hello, World!"
}
```

#### Broadcast Message

```json
{
  "type": "broadcast",
  "sender": "your username",
  "content": "Hello, World!"
}
```

#### Welcome Message

```json
{
  "type": "welcome",
  "client_id": "unique_id",
  "your_username": "username",
  "client_list": ["user1", "user2", "user3"]
}
```

## Usage Examples

### Connecting to the Server

```javascript
// Browser JavaScript example
const ws = new WebSocket("ws://localhost:8080/ws");
```

### Sending a Broadcast Message

```javascript
ws.send(JSON.stringify({
    type: "broadcast",
    sender: "username"
    content: "Hello everyone!"
}));
```

### Sending a Private Message

```javascript
ws.send(
  JSON.stringify({
    type: "private",
    receiver: "<ID>",
    content: "Hi Jane!",
  })
);
```

## Error Handling

The server handles various error scenarios:

- Invalid JSON messages
- Connection timeouts
- Client disconnections

## Project Structure

```
.
├── server
│   ├── client.go
│   ├── handlers.go
│   ├── server.go
│   ├── types.go
├── go.mod
├── go.sum
├── main.go
├── README.md
```

## Details

- Ayush Kumar Awasthi
- [Email](mailto:ayushawasthi2409.gmail.com)
- Phone Number : 9756798580

## Acknowledgments

- [Gorilla WebSocket](https://github.com/gorilla/websocket) for the WebSocket implementation
- The Go community for their excellent documentation and support
