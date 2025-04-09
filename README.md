# HTTP Server Implementation

A lightweight HTTP/1.1 server built from scratch using Go.

## Features

- HTTP/1.1 compliant server implementation  
- Supports GET requests  
- Proper header parsing and generation  
- TCP socket connection management  
- Error handling with appropriate status codes

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/http-server.git
cd http-server

# Build the server
go run ./cmd/httpserver/main.go
```
## Usage
You can use any HTTP client to make requests. Here's an example using Node.js:
Using Node.js http module:
```bash
const http = require('http');

http.get('http://localhost:8080/hello', (res) => {
  let data = '';

  res.on('data', chunk => {
    data += chunk;
  });

  res.on('end', () => {
    console.log('Response:', data);
  });
}).on('error', err => {
  console.error('Error:', err.message);
});
```
